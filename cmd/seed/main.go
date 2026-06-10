package main

import (
	"fmt"
	"log"

	"kimsha/internal/config"
	"kimsha/internal/models"
	"kimsha/internal/repository"
	"kimsha/pkg/password"

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

	if cfg.Admin.Email == "" || cfg.Admin.Password == "" {
		log.Fatal("ADMIN_EMAIL and ADMIN_PASSWORD must be set in .env")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Africa/Addis_Ababa",
		cfg.Postgres.Host, cfg.Postgres.Port,
		cfg.Postgres.User, cfg.Postgres.Password,
		cfg.Postgres.DBName, cfg.Postgres.SSL,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("db:", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Tenant{}); err != nil {
		log.Fatal("migrate:", err)
	}

	adminRepo := repository.NewAdminRepository(db)

	exists, err := adminRepo.SuperAdminExists()
	if err != nil {
		log.Fatal("check:", err)
	}
	if exists {
		log.Println("super admin already exists — skipping")
		return
	}

	hashed, err := password.Hash(cfg.Admin.Password)
	if err != nil {
		log.Fatal("hash:", err)
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

	if err := adminRepo.CreateSuperAdmin(admin); err != nil {
		log.Fatal("create:", err)
	}

	log.Printf("super admin created: %s", cfg.Admin.Email)
}
