package handlers

import (
	"math"
	"strconv"
	"strings"
	"transaction-service/dto"
	"transaction-service/services"

	"github.com/go-playground/validator/v10"
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
			Message: "Invalid request body: " + err.Error(),
		})
	}

	if len(req.Items) == 0 {
		return c.Status(400).JSON(dto.ApiResponse{
			Success: false,
			Message: "Transaction items are required",
		})
	}
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		errs := err.(validator.ValidationErrors)
		var messages []string
		for _, e := range errs {
			switch e.Tag() {
			case "required":
				messages = append(messages, e.Field()+" is required")
			case "gt":
				messages = append(messages, e.Field()+" must be greater than "+e.Param())
			}
		}

		return c.Status(400).JSON(dto.ApiResponse{
			Success: false,
			Message: strings.Join(messages, ", "),
		})
	}

	transaction, err := h.service.CreateTransaction(&req)
	if err != nil {
		statusCode := 500
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "insufficient stock") {
			statusCode = 400
		}

		return c.Status(statusCode).JSON(dto.ApiResponse{
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
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	sortBy := c.Query("sort_by", "created_at")
	order := strings.ToUpper(c.Query("order", "DESC"))

	// validasai param page dan limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	transactions, total, err := h.service.GetAllTransactions(page, limit, search, sortBy, order)
	if err != nil {
		return c.Status(500).JSON(dto.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.JSON(dto.ApiResponse{
		Success: true,
		Message: "Transactions retrieved successfully",
		Data: map[string]interface{}{
			"transactions": transactions,
			"pagination": map[string]interface{}{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": totalPages,
				"has_next":    page < totalPages,
				"has_prev":    page > 1,
			},
		},
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
		statusCode := 500
		if strings.Contains(err.Error(), "not found") {
			statusCode = 404
		}

		return c.Status(statusCode).JSON(dto.ApiResponse{
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
