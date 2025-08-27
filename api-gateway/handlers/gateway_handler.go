package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

type GatewayHandler struct{}

func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{}
}

func (h *GatewayHandler) ProductProxy(c *fiber.Ctx) error {
	return proxy.Do(c, "http://product-service:8081"+c.OriginalURL())
}

func (h *GatewayHandler) TransactionProxy(c *fiber.Ctx) error {
	return proxy.Do(c, "http://transaction-service:8082"+c.OriginalURL())
}