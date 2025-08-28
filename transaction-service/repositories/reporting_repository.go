package repositories

import (
	"fmt"
	"transaction-service/config"
	"transaction-service/dto"
)

type ReportingRepository interface {
	GetTransactionSummary(filter dto.ReportingFilterDTO) ([]dto.TransactionSummaryDTO, error)
	GetProductSalesReport(filter dto.ReportingFilterDTO) ([]dto.ProductSalesReportDTO, error)
	GetLowStockAlert() ([]dto.LowStockAlertDTO, error)
}

type reportingRepository struct{}

func NewReportingRepository() ReportingRepository {
	return &reportingRepository{}
}

func (r *reportingRepository) GetTransactionSummary(filter dto.ReportingFilterDTO) ([]dto.TransactionSummaryDTO, error) {
	query := `
		SELECT 
			id,
			transaction_date,
			total_amount,
			total_items,
			total_quantity
		FROM v_transaction_summary
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	// filter berdasarkan tanggal
	if filter.StartDate != nil {
		query += " AND transaction_date >= $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		query += " AND transaction_date <= $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	// Add pagination
	if filter.Limit > 0 {
		query += " LIMIT $" + fmt.Sprintf("%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += " OFFSET $" + fmt.Sprintf("%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []dto.TransactionSummaryDTO
	for rows.Next() {
		var summary dto.TransactionSummaryDTO
		err := rows.Scan(
			&summary.ID,
			&summary.TransactionDate,
			&summary.TotalAmount,
			&summary.TotalItems,
			&summary.TotalQuantity,
		)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (r *reportingRepository) GetProductSalesReport(filter dto.ReportingFilterDTO) ([]dto.ProductSalesReportDTO, error) {
	query := `
		SELECT 
			id,
			product_name,
			current_price,
			current_stock,
			total_sold,
			total_revenue
		FROM v_product_sales_report
		ORDER BY total_sold DESC
	`

	args := []interface{}{}

	// Add pagination
	if filter.Limit > 0 {
		query += " LIMIT $1"
		args = append(args, filter.Limit)

		if filter.Offset > 0 {
			query += " OFFSET $2"
			args = append(args, filter.Offset)
		}
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []dto.ProductSalesReportDTO
	for rows.Next() {
		var report dto.ProductSalesReportDTO
		err := rows.Scan(
			&report.ID,
			&report.ProductName,
			&report.CurrentPrice,
			&report.CurrentStock,
			&report.TotalSold,
			&report.TotalRevenue,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func (r *reportingRepository) GetLowStockAlert() ([]dto.LowStockAlertDTO, error) {
	query := `
		SELECT  id, name, price, stock, stock_status
		FROM v_low_stock_alert
		ORDER BY stock ASC
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []dto.LowStockAlertDTO
	for rows.Next() {
		var alert dto.LowStockAlertDTO
		err := rows.Scan(
			&alert.ID,
			&alert.Name,
			&alert.Price,
			&alert.Stock,
			&alert.StockStatus,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}
