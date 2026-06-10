package routes

import (
	"kimsha/internal/api/handlers"
	"kimsha/internal/api/middleware"
	"kimsha/internal/models"

	"github.com/gofiber/fiber/v2"
	gows "github.com/gofiber/websocket/v2"
)

type Handlers struct {
	Auth    *handlers.AuthHandler
	Admin   *handlers.AdminHandler
	Menu    *handlers.MenuHandler
	Table   *handlers.TableHandler
	Order   *handlers.OrderHandler
	Kitchen *handlers.KitchenHandler
	Payment *handlers.PaymentHandler
	Report  *handlers.ReportHandler
	Sync    *handlers.SyncHandler
	User    *handlers.UserHandler
	Tenant  *handlers.TenantHandler
	Upload  *handlers.UploadHandler
}

func Register(app *fiber.App, h *Handlers) {
	api := app.Group("/api/v1")

	// Super-admin — public login, then protected endpoints
	api.Post("/admin/login", h.Admin.Login)
	adminGroup := api.Group("/admin", middleware.Auth(), middleware.RequireRole(string(models.RoleSuperAdmin)))
	adminGroup.Get("/stats", h.Admin.Stats)
	adminGroup.Get("/tenants", h.Admin.ListTenants)
	adminGroup.Post("/tenants", h.Admin.CreateTenant)
	adminGroup.Get("/tenants/:id", h.Admin.GetTenant)
	adminGroup.Patch("/tenants/:id/status", h.Admin.SetTenantStatus)

	// Public
	auth := api.Group("/auth")
	auth.Post("/register", h.Auth.Register)
	auth.Post("/login", h.Auth.Login)
	auth.Post("/pin-login", h.Auth.PINLogin)

	// Protected
	protected := api.Use(middleware.Auth())

	// Tenant
	protected.Get("/tenant", h.Tenant.Get)
	protected.Put("/tenant", middleware.RequireRole(string(models.RoleOwner), string(models.RoleManager)), h.Tenant.Update)

	// Users — owner/manager only
	staff := protected.Group("/users", middleware.RequireRole(string(models.RoleOwner), string(models.RoleManager)))
	staff.Get("/", h.User.List)
	staff.Post("/", h.User.Create)
	staff.Get("/:id", h.User.Get)
	staff.Put("/:id", h.User.Update)
	staff.Delete("/:id", h.User.Delete)
	staff.Put("/:id/pin", h.User.SetPIN)

	// Menu
	menu := protected.Group("/menu")
	menu.Get("/categories", h.Menu.ListCategories)
	menu.Post("/categories", middleware.RequireRole(string(models.RoleOwner), string(models.RoleManager)), h.Menu.CreateCategory)
	menu.Put("/categories/:id", middleware.RequireRole(string(models.RoleOwner), string(models.RoleManager)), h.Menu.UpdateCategory)
	menu.Delete("/categories/:id", middleware.RequireRole(string(models.RoleOwner), string(models.RoleManager)), h.Menu.DeleteCategory)

	menu.Get("/items", h.Menu.ListItems)
	menu.Post("/items", middleware.RequireRole(string(models.RoleOwner), string(models.RoleManager)), h.Menu.CreateItem)
	menu.Get("/items/:id", h.Menu.GetItem)
	menu.Put("/items/:id", middleware.RequireRole(string(models.RoleOwner), string(models.RoleManager)), h.Menu.UpdateItem)
	menu.Delete("/items/:id", middleware.RequireRole(string(models.RoleOwner), string(models.RoleManager)), h.Menu.DeleteItem)
	menu.Put("/items/:id/toggle", middleware.RequireRole(string(models.RoleOwner), string(models.RoleManager)), h.Menu.ToggleItem)

	// Tables
	tables := protected.Group("/tables")
	tables.Get("/", h.Table.List)
	tables.Post("/", middleware.RequireRole(string(models.RoleOwner), string(models.RoleManager)), h.Table.Create)
	tables.Put("/:id", middleware.RequireRole(string(models.RoleOwner), string(models.RoleManager)), h.Table.Update)
	tables.Delete("/:id", middleware.RequireRole(string(models.RoleOwner), string(models.RoleManager)), h.Table.Delete)
	tables.Put("/:id/status", h.Table.UpdateStatus)

	// Orders
	orders := protected.Group("/orders")
	orders.Get("/", h.Order.List)
	orders.Post("/", h.Order.Create)
	orders.Get("/:id", h.Order.Get)
	orders.Put("/:id/status", h.Order.UpdateStatus)
	orders.Post("/:id/items", h.Order.AddItem)
	orders.Put("/:id/items/:iid/status", h.Order.UpdateItemStatus)
	orders.Delete("/:id/items/:iid", h.Order.RemoveItem)

	// Kitchen
	kitchen := protected.Group("/kitchen")
	kitchen.Get("/tickets", h.Kitchen.ActiveTickets)
	kitchen.Put("/tickets/:id/status", h.Kitchen.UpdateTicketStatus)

	// WebSocket upgrade — use middleware.Auth before upgrading
	app.Use("/api/v1/kitchen/ws", func(c *fiber.Ctx) error {
		if gows.IsWebSocketUpgrade(c) {
			return middleware.Auth()(c)
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/api/v1/kitchen/ws", gows.New(h.Kitchen.WebSocket))

	// Payments
	protected.Post("/orders/:id/pay", h.Payment.Pay)
	protected.Get("/payments", h.Payment.List)
	protected.Post("/cash/open", middleware.RequireRole(string(models.RoleCashier), string(models.RoleOwner)), h.Payment.OpenShift)
	protected.Post("/cash/close", middleware.RequireRole(string(models.RoleCashier), string(models.RoleOwner)), h.Payment.CloseShift)
	protected.Post("/cash/in", middleware.RequireRole(string(models.RoleCashier), string(models.RoleOwner)), h.Payment.CashIn)
	protected.Post("/cash/out", middleware.RequireRole(string(models.RoleCashier), string(models.RoleOwner)), h.Payment.CashOut)
	protected.Get("/cash/summary", h.Payment.CashSummary)

	// Reports
	reports := protected.Group("/reports", middleware.RequireRole(string(models.RoleOwner), string(models.RoleManager)))
	reports.Get("/daily", h.Report.Daily)
	reports.Get("/items", h.Report.TopItems)
	reports.Get("/hourly", h.Report.HourlySales)
	reports.Get("/waiters", h.Report.WaiterStats)

	// Sync
	sync := protected.Group("/sync")
	sync.Post("/push", h.Sync.Push)
	sync.Get("/menu", h.Sync.MenuSnapshot)
	sync.Get("/tables", h.Sync.TableSnapshot)

	// Upload
	protected.Post("/upload/image", h.Upload.UploadImage)
}
