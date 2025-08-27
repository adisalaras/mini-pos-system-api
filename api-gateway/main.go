package main

import (
	"api-gateway/handlers"
	"api-gateway/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Middleware
	app.Use(cors.New())
	app.Use(handlers.LoggingMiddleware)

	// Routes
	routes.SetupRoutes(app)

	log.Println("API Gateway running on :8080")
	log.Fatal(app.Listen(":8080"))
}
