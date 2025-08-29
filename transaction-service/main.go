package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"transaction-service/config"
	"transaction-service/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if exists (for development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	config.InitDatabase()
	defer config.CloseDatabase()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			log.Printf("Error: %v", err)
			return ctx.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} - ${latency}\n",
	}))
	app.Use(recover.New())

	routes.SetupTransactionRoutes(app)

	app.Get("/health/db", func(c *fiber.Ctx) error {
		err := config.PingDatabase()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": "error",
				"error":  err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "connected",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Transaction Service running on :%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Panic(err)
		}
	}()

	<-c
	log.Println("Gracefully shutting down")
	app.Shutdown()
	log.Println("Transaction service shutdown complete")
}
