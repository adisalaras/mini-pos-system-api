package dto

type CreateProductRequest struct {
	Name  string  `json:"name" validate:"required,min=1,max=100"`
	Price float64 `json:"price" validate:"required,min=0"`
	Stock int     `json:"stock" validate:"required,min=0"`
}

type UpdateProductRequest struct {
	Name  string  `json:"name,omitempty"`
	Price float64 `json:"price,omitempty"`
	Stock int     `json:"stock,omitempty"`
}

type ProductResponse struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}