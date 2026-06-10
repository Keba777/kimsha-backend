package main

import (
	"fmt"
	"log"

	"kimsha/internal/api/handlers"
	"kimsha/internal/api/middleware"
	"kimsha/internal/api/routes"
	"kimsha/internal/config"
	"kimsha/internal/models"
	"kimsha/internal/repository"
	"kimsha/internal/services"
	"kimsha/pkg/jwt"
	"kimsha/pkg/password"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config:", err)
	}

	// Init JWT
	jwt.Init(cfg.JWT.Secret)

	// Connect PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Africa/Addis_Ababa",
		cfg.Postgres.Host, cfg.Postgres.Port,
		cfg.Postgres.User, cfg.Postgres.Password,
		cfg.Postgres.DBName, cfg.Postgres.SSL,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("db:", err)
	}

	// AutoMigrate
	if err := db.AutoMigrate(
		&models.Tenant{},
		&models.User{},
		&models.Category{},
		&models.MenuItem{},
		&models.AddOn{},
		&models.Table{},
		&models.Order{},
		&models.OrderItem{},
		&models.KitchenTicket{},
		&models.Payment{},
		&models.CashTransaction{},
		&models.SyncLog{},
		&models.AuditLog{},
	); err != nil {
		log.Fatal("migrate:", err)
	}

	// Repositories
	authRepo := repository.NewAuthRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	// Seed super admin on first boot (no-op if already exists)
	if err := seedSuperAdmin(cfg, adminRepo); err != nil {
		log.Printf("warn: super admin seed failed: %v", err)
	}
	menuRepo := repository.NewMenuRepository(db)
	tableRepo := repository.NewTableRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	kitchenRepo := repository.NewKitchenRepository(db)
	payRepo := repository.NewPaymentRepository(db)
	reportRepo := repository.NewReportRepository(db)
	syncRepo := repository.NewSyncRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Storage (MinIO)
	storageSvc, err := services.NewStorageService(cfg.MinIO)
	if err != nil {
		log.Fatal("minio:", err)
	}

	// Services
	authSvc := services.NewAuthService(authRepo, cfg.JWT.Expiry)
	menuSvc := services.NewMenuService(menuRepo)
	tableSvc := services.NewTableService(tableRepo)
	orderSvc := services.NewOrderService(orderRepo, menuRepo, kitchenRepo, tableRepo)
	kitchenSvc := services.NewKitchenService(kitchenRepo)
	paymentSvc := services.NewPaymentService(payRepo, orderRepo, tableRepo)
	reportSvc := services.NewReportService(reportRepo, orderRepo)
	syncSvc := services.NewSyncService(syncRepo, orderRepo, menuRepo, tableRepo)

	// Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "ቅምሻ API",
		ErrorHandler: errorHandler,
	})

	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS(cfg.CORS.Origins))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "app": "ቅምሻ"})
	})

	// Register all routes
	routes.Register(app, &routes.Handlers{
		Auth:    handlers.NewAuthHandler(authSvc, cfg.Admin.AllowSelfRegister),
		Admin:   handlers.NewAdminHandler(adminRepo, cfg.JWT.Expiry),
		Menu:    handlers.NewMenuHandler(menuSvc),
		Table:   handlers.NewTableHandler(tableSvc),
		Order:   handlers.NewOrderHandler(orderSvc),
		Kitchen: handlers.NewKitchenHandler(kitchenSvc),
		Payment: handlers.NewPaymentHandler(paymentSvc),
		Report:  handlers.NewReportHandler(reportSvc),
		Sync:    handlers.NewSyncHandler(syncSvc),
		User:    handlers.NewUserHandler(userRepo),
		Tenant:  handlers.NewTenantHandler(db),
		Upload:  handlers.NewUploadHandler(storageSvc),
	})

	log.Printf("ቅምሻ server starting on :%s", cfg.App.Port)
	log.Fatal(app.Listen(":" + cfg.App.Port))
}

func seedSuperAdmin(cfg *config.Config, repo *repository.AdminRepository) error {
	if cfg.Admin.Email == "" || cfg.Admin.Password == "" {
		return nil
	}
	exists, err := repo.SuperAdminExists()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	hashed, err := password.Hash(cfg.Admin.Password)
	if err != nil {
		return err
	}
	admin := &models.User{
		ID:       uuid.New(),
		TenantID: uuid.Nil,
		Name:     "Super Admin",
		Email:    cfg.Admin.Email,
		Password: hashed,
		Role:     models.RoleSuperAdmin,
		IsActive: true,
	}
	if err := repo.CreateSuperAdmin(admin); err != nil {
		return err
	}
	log.Printf("super admin created: %s", cfg.Admin.Email)
	return nil
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{"success": false, "error": err.Error()})
}
