package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

type GatewayHandler struct{}

func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{}
}

func (h *GatewayHandler) ProductProxy(c *fiber.Ctx) error {
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://127.0.0.1:8081"
	}
	return proxy.Do(c, productServiceURL+c.OriginalURL())
}

func (h *GatewayHandler) TransactionProxy(c *fiber.Ctx) error {
	transactionServiceURL := os.Getenv("TRANSACTION_SERVICE_URL")
	if transactionServiceURL == "" {
		transactionServiceURL = "http://127.0.0.1:8082"
	}
	return proxy.Do(c, transactionServiceURL+c.OriginalURL())
}
func (h *GatewayHandler) ReportingProxy(c *fiber.Ctx) error {
	reportingServiceURL := os.Getenv("REPORTING_SERVICE_URL")
	if reportingServiceURL == "" {
		reportingServiceURL = "http://127.0.0.1:8082"
	}
	return proxy.Do(c, reportingServiceURL+c.OriginalURL())
}
