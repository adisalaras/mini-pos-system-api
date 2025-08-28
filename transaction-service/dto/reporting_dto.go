package dto

import "time"

//summary transaction
type TransactionSummaryDTO struct {
	ID              uint      `json:"id"`
	TransactionDate time.Time `json:"transaction_date"`
	TotalAmount     float64   `json:"total_amount"`
	TotalItems      int       `json:"total_items"`
	TotalQuantity   int       `json:"total_quantity"`
}

// sales report per product
type ProductSalesReportDTO struct {
	ID           uint    `json:"id"`
	ProductName  string  `json:"product_name"`
	CurrentPrice float64 `json:"current_price"`
	CurrentStock int     `json:"current_stock"`
	TotalSold    int     `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}

// alert jika stock produk menipis
type LowStockAlertDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	StockStatus string  `json:"stock_status"`
}

// filter untuk laporan
type ReportingFilterDTO struct {
	StartDate *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" form:"end_date"`
	Limit     int        `json:"limit,omitempty" form:"limit"`
	Offset    int        `json:"offset,omitempty" form:"offset"`
}
