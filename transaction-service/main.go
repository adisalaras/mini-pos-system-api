package main

import (
	"log"
	"os"
	"transaction-service/config"
	"transaction-service/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize database
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	config.InitDatabase()
	defer config.CloseDatabase()

	app := fiber.New()
	app.Use(cors.New())

	// Setup routes
	routes.SetupTransactionRoutes(app)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Transaction Service is running")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Transaction Service running on :%s\n", port)
	log.Fatal(app.Listen(":" + port))

}
