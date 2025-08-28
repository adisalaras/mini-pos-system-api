package handlers

import (
	"product-service/dto"
	"product-service/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ApiResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// validasii 
	if req.Name == "" {
		return c.Status(400).JSON(dto.ApiResponse{
			Success: false,
			Message: "Product name is required",
		})
	}
	if req.Price <= 0 {
		return c.Status(400).JSON(dto.ApiResponse{
			Success: false,
			Message: "Price must be greater than 0",
		})
	}
	if req.Stock < 0 {
		return c.Status(400).JSON(dto.ApiResponse{
			Success: false,
			Message: "Stock cannot be negative",
		})
	}

	product, err := h.service.CreateProduct(&req)
	if err != nil {
		return c.Status(500).JSON(dto.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(201).JSON(dto.ApiResponse{
		Success: true,
		Message: "Product created successfully",
		Data:    product,
	})
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.service.GetAllProducts()
	if err != nil {
		return c.Status(500).JSON(dto.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(dto.ApiResponse{
		Success: true,
		Message: "Products retrieved successfully",
		Data:    products,
	})
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(400).JSON(dto.ApiResponse{
			Success: false,
			Message: "Invalid product ID",
		})
	}

	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(dto.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(dto.ApiResponse{
		Success: true,
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(400).JSON(dto.ApiResponse{
			Success: false,
			Message: "Invalid product ID",
		})
	}

	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ApiResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	product, err := h.service.UpdateProduct(uint(id), &req)
	if err != nil {
		return c.Status(500).JSON(dto.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(dto.ApiResponse{
		Success: true,
		Message: "Product updated successfully",
		Data:    product,
	})
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(400).JSON(dto.ApiResponse{
			Success: false,
			Message: "Invalid product ID",
		})
	}

	err = h.service.DeleteProduct(uint(id))
	if err != nil {
		return c.Status(500).JSON(dto.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(dto.ApiResponse{
		Success: true,
		Message: "Product deleted successfully",
	})
}
