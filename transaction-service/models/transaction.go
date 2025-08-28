package models

import (
	"time"
)

type Transaction struct {
	ID               uint              `json:"id"`
	TransactionDate  time.Time         `json:"transaction_date"`
	TotalAmount      float64           `json:"total_amount"`
	TransactionItems []TransactionItem `json:"transaction_items"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	DeletedAt        *time.Time        `json:"deleted_at,omitempty"`
}

type TransactionItem struct {
	ID            uint      `json:"id"`
	TransactionID uint      `json:"transaction_id"`
	ProductID     uint      `json:"product_id"`
	ProductName   string    `json:"product_name"`
	Price         float64   `json:"price"`
	Quantity      int       `json:"quantity"`
	Subtotal      float64   `json:"subtotal"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}


