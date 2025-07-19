package routes

import (
	"pos-bengkel-backend/controllers"
	"pos-bengkel-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Initialize controllers
	authController := controllers.NewAuthController()
	vehicleController := controllers.NewVehicleController()
	customerController := controllers.NewCustomerController()
	transactionController := controllers.NewTransactionController()

	// API group
	api := app.Group("/api")

	// Authentication routes (public)
	auth := api.Group("/auth")
	auth.Post("/login", authController.Login)
	auth.Post("/register", authController.Register)
	auth.Get("/profile", middleware.JWTMiddleware(), authController.GetProfile)
	auth.Post("/refresh", middleware.JWTMiddleware(), authController.RefreshToken)

	// Vehicle routes (public read, authenticated write)
	vehicles := api.Group("/vehicles")
	vehicles.Get("/", middleware.OptionalJWTMiddleware(), vehicleController.GetVehicles)
	vehicles.Get("/search", middleware.OptionalJWTMiddleware(), vehicleController.SearchVehicles)
	vehicles.Get("/:id", middleware.OptionalJWTMiddleware(), vehicleController.GetVehicle)
	vehicles.Post("/", middleware.JWTMiddleware(), middleware.RequireRole("admin", "sales"), vehicleController.CreateVehicle)
	vehicles.Put("/:id", middleware.JWTMiddleware(), middleware.RequireRole("admin", "sales"), vehicleController.UpdateVehicle)
	vehicles.Delete("/:id", middleware.JWTMiddleware(), middleware.RequireRole("admin"), vehicleController.DeleteVehicle)

	// Customer routes (authenticated)
	customers := api.Group("/customers")
	customers.Use(middleware.JWTMiddleware())
	customers.Get("/", middleware.RequireRole("admin", "sales", "cashier"), customerController.GetCustomers)
	customers.Get("/:id", middleware.RequireRole("admin", "sales", "cashier"), customerController.GetCustomer)
	customers.Post("/", middleware.RequireRole("admin", "sales"), customerController.CreateCustomer)
	customers.Put("/:id", middleware.RequireRole("admin", "sales"), customerController.UpdateCustomer)
	customers.Delete("/:id", middleware.RequireRole("admin"), customerController.DeleteCustomer)

	// Transaction routes (authenticated)
	transactions := api.Group("/transactions")
	transactions.Use(middleware.JWTMiddleware())
	transactions.Get("/", middleware.RequireRole("admin", "sales", "cashier"), transactionController.GetTransactions)
	transactions.Get("/:id", middleware.RequireRole("admin", "sales", "cashier"), transactionController.GetTransaction)
	transactions.Post("/", middleware.RequireRole("admin", "sales", "cashier"), transactionController.CreateTransaction)
	transactions.Put("/:id/status", middleware.RequireRole("admin", "sales", "cashier"), transactionController.UpdateTransactionStatus)

	// Admin-only user management routes
	users := api.Group("/users")
	users.Use(middleware.JWTMiddleware(), middleware.RequireRole("admin"))
	// Note: User management endpoints can be added here in the future
	// users.Get("/", userController.GetUsers)
	// users.Post("/", userController.CreateUser)
	// users.Get("/:id", userController.GetUser)
	// users.Put("/:id", userController.UpdateUser)
	// users.Delete("/:id", userController.DeleteUser)

	// Test drive routes (future implementation)
	testDrives := api.Group("/test-drives")
	testDrives.Use(middleware.JWTMiddleware())
	// Note: Test drive endpoints can be added here in the future
	// testDrives.Get("/", testDriveController.GetTestDrives)
	// testDrives.Post("/", testDriveController.CreateTestDrive)
	// testDrives.Get("/:id", testDriveController.GetTestDrive)
	// testDrives.Put("/:id/status", testDriveController.UpdateTestDriveStatus)
}