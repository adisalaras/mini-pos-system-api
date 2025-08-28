package models

import (
	"time"
)

type Product struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Price     float64    `json:"price"`
	Stock     int        `json:"stock"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}