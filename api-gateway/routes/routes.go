package routes

import (
	"api-gateway/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	gatewayHandler := handlers.NewGatewayHandler()

	// Product service routes
	products := app.Group("/api/products")
	products.Use(gatewayHandler.ProductProxy)

	// Transaction service routes
	transactions := app.Group("/api/transactions")
	transactions.Use(gatewayHandler.TransactionProxy)

	// Reporting service routes
	reports := app.Group("/api/reports")
	reports.Use(gatewayHandler.ReportingProxy)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
