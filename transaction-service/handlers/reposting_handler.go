package handlers

import (
	"strconv"
	"time"
	"transaction-service/dto"
	"transaction-service/services"

	"github.com/gofiber/fiber/v2"
)

type ReportingHandler struct {
	reportingService services.ReportingService
}

func NewReportingHandler(reportingService services.ReportingService) *ReportingHandler {
	return &ReportingHandler{reportingService: reportingService}
}

func (h *ReportingHandler) GetTransactionSummary(c *fiber.Ctx) error {
	filter := dto.ReportingFilterDTO{}

	// jika start_date diset
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}

	// jika end_date diset
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}

	// jika limit diset
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	// jika offset diset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}

	summaries, err := h.reportingService.GetTransactionSummary(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get transaction summary",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Transaction summary retrieved successfully",
		"data":    summaries,
		"count":   len(summaries),
	})
}

func (h *ReportingHandler) GetProductSalesReport(c *fiber.Ctx) error {
	// memisahkan parameter filter dari query
	filter := dto.ReportingFilterDTO{}

	// jika limit diset
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	// jika offset diset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}

	reports, err := h.reportingService.GetProductSalesReport(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get product sales report",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Product sales report retrieved successfully",
		"data":    reports,
		"count":   len(reports),
	})
}

func (h *ReportingHandler) GetLowStockAlert(c *fiber.Ctx) error {
	alerts, err := h.reportingService.GetLowStockAlert()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get low stock alerts",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Low stock alerts retrieved successfully",
		"data":    alerts,
		"count":   len(alerts),
	})
}

func (h *ReportingHandler) GetDashboardSummary(c *fiber.Ctx) error {
	dashboard, err := h.reportingService.GetDashboardSummary()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get dashboard summary",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Dashboard summary retrieved successfully",
		"data":    dashboard,
	})
}
