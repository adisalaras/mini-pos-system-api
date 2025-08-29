package routes

import (
	"os"
	"transaction-service/clients"
	"transaction-service/handlers"
	"transaction-service/repositories"
	"transaction-service/services"

	"github.com/gofiber/fiber/v2"
)

func SetupTransactionRoutes(app *fiber.App) {
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://product-service:8081" 
	}

	productClient := clients.NewProductClient(productServiceURL)
	transactionRepo := repositories.NewTransactionRepository()
	transactionService := services.NewTransactionService(transactionRepo, productClient)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	reportingRepo := repositories.NewReportingRepository()
	reportingService := services.NewReportingService(reportingRepo)
	reportingHandler := handlers.NewReportingHandler(reportingService)

	api := app.Group("/api")

	transactions := api.Group("/transactions")
	transactions.Post("/", transactionHandler.CreateTransaction)
	transactions.Get("/", transactionHandler.GetAllTransactions)
	transactions.Get("/:id", transactionHandler.GetTransaction)

	reports := api.Group("/reports")
	reports.Get("/transactions", reportingHandler.GetTransactionSummary)
	reports.Get("/products", reportingHandler.GetProductSalesReport)
	reports.Get("/low-stock", reportingHandler.GetLowStockAlert)
	reports.Get("/dashboard", reportingHandler.GetDashboardSummary)
}
