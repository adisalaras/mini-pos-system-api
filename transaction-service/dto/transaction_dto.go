package dto

type CreateTransactionRequest struct {
	Items []TransactionItemRequest `json:"items" validate:"required,min=1"`
}

type TransactionItemRequest struct {
	ProductID   uint    `json:"product_id" validate:"required"`
	ProductName string  `json:"product_name" validate:"required"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
}

type TransactionResponse struct {
	ID               uint                      `json:"id"`
	TransactionDate  string                    `json:"transaction_date"`
	TotalAmount      float64                   `json:"total_amount"`
	TransactionItems []TransactionItemResponse `json:"transaction_items"`
	CreatedAt        string                    `json:"created_at"`
}

type TransactionItemResponse struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

