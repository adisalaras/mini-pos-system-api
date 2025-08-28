package routes

import (
	"transaction-service/handlers"
	"transaction-service/repositories"
	"transaction-service/services"
	"github.com/gofiber/fiber/v2"
)

func SetupTransactionRoutes(app *fiber.App) {
	// Initialize dependencies
	transactionRepo := repositories.NewTransactionRepository()
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	
	// Routes
	api := app.Group("/api")
	transactions := api.Group("/transactions")
	
	transactions.Post("/", transactionHandler.CreateTransaction)
	transactions.Get("/", transactionHandler.GetAllTransactions)
	transactions.Get("/:id", transactionHandler.GetTransaction)
}