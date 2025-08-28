package handlers

import (
	"strconv"
	"transaction-service/dto"
	"transaction-service/services"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var req dto.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ApiResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if len(req.Items) == 0 {
		return c.Status(400).JSON(dto.ApiResponse{
			Success: false,
			Message: "Transaction items are required",
		})
	}

	transaction, err := h.service.CreateTransaction(&req)
	if err != nil {
		return c.Status(500).JSON(dto.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(201).JSON(dto.ApiResponse{
		Success: true,
		Message: "Transaction created successfully",
		Data:    transaction,
	})
}

func (h *TransactionHandler) GetAllTransactions(c *fiber.Ctx) error {
	transactions, err := h.service.GetAllTransactions()
	if err != nil {
		return c.Status(500).JSON(dto.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(dto.ApiResponse{
		Success: true,
		Message: "Transactions retrieved successfully",
		Data:    transactions,
	})
}

func (h *TransactionHandler) GetTransaction(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(400).JSON(dto.ApiResponse{
			Success: false,
			Message: "Invalid transaction ID",
		})
	}

	transaction, err := h.service.GetTransactionByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(dto.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(dto.ApiResponse{
		Success: true,
		Message: "Transaction retrieved successfully",
		Data:    transaction,
	})
}
