package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggingMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	err := c.Next()

	log.Printf("%s %s - %d - %v", c.Method(), c.Path(), c.Response().StatusCode(), time.Since(start))

	return err
}
