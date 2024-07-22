package router

import (
	"maglo/handler"
	"maglo/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Get("/google", handler.Auth)
	auth.Get("/google/callback/", handler.Callback)
	auth.Post("/logout", middleware.Protected(), handler.Logout)

	// User
	user := api.Group("/user")
	user.Get("/:id", handler.GetUser)
	user.Post("/", handler.CreateUser)
	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	// Transaction
	product := api.Group("/transaction")
	product.Get("/", handler.GetAllTransactions)
	product.Get("/:id", handler.GetTransaction)
	product.Post("/", middleware.Protected(), handler.CreateTransaction)
	product.Delete("/:id", middleware.Protected(), handler.DeleteTransaction)
}