package services

import (
	"transaction-service/dto"
	"transaction-service/repositories"
)

type ReportingService interface {
	GetTransactionSummary(filter dto.ReportingFilterDTO) ([]dto.TransactionSummaryDTO, error)
	GetProductSalesReport(filter dto.ReportingFilterDTO) ([]dto.ProductSalesReportDTO, error)
	GetLowStockAlert() ([]dto.LowStockAlertDTO, error)
	GetDashboardSummary() (map[string]interface{}, error)
}

type reportingService struct {
	reportingRepo repositories.ReportingRepository
}

func NewReportingService(reportingRepo repositories.ReportingRepository) ReportingService {
	return &reportingService{reportingRepo: reportingRepo}
}

func (s *reportingService) GetTransactionSummary(filter dto.ReportingFilterDTO) ([]dto.TransactionSummaryDTO, error) {
	// jika limit tidak diset, set default 50
	if filter.Limit == 0 {
		filter.Limit = 50
	}

	return s.reportingRepo.GetTransactionSummary(filter)
}

func (s *reportingService) GetProductSalesReport(filter dto.ReportingFilterDTO) ([]dto.ProductSalesReportDTO, error) {
	// jika limit tidak diset, set default 100
	if filter.Limit == 0 {
		filter.Limit = 100
	}

	return s.reportingRepo.GetProductSalesReport(filter)
}

func (s *reportingService) GetLowStockAlert() ([]dto.LowStockAlertDTO, error) {
	return s.reportingRepo.GetLowStockAlert()
}

func (s *reportingService) GetDashboardSummary() (map[string]interface{}, error) {
	// mendapatkan 10 transaksi terbaru
	recentFilter := dto.ReportingFilterDTO{Limit: 10}
	recentTransactions, err := s.reportingRepo.GetTransactionSummary(recentFilter)
	if err != nil {
		return nil, err
	}

	// mendapatkan 5 produk teratas
	topProductsFilter := dto.ReportingFilterDTO{Limit: 5}
	topProducts, err := s.reportingRepo.GetProductSalesReport(topProductsFilter)
	if err != nil {
		return nil, err
	}

	// mendapatkan alert stock menipis
	lowStockAlerts, err := s.reportingRepo.GetLowStockAlert()
	if err != nil {
		return nil, err
	}

	// hitung total transaksi dan total revenue
	var totalRevenue float64
	var totalTransactions int
	for _, transaction := range recentTransactions {
		totalRevenue += transaction.TotalAmount
		totalTransactions++
	}

	dashboard := map[string]interface{}{
		"total_transactions":  totalTransactions,
		"total_revenue":       totalRevenue,
		"recent_transactions": recentTransactions,
		"top_products":        topProducts,
		"low_stock_alerts":    lowStockAlerts,
		"low_stock_count":     len(lowStockAlerts),
	}

	return dashboard, nil
}
