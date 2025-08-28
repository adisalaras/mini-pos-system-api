package main

import (
	"log"
	"os"
	"product-service/config"

	// "product-service/routes"

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
	// routes.SetupProductRoutes(app)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("✅ Product Service is running")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Product Service running on :%s\n", port)
	log.Fatal(app.Listen(":" + port))



}
