package repositories

import (
	"database/sql"
	"fmt"
	"product-service/config"
	"product-service/models"
	"time"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetAll(page, limit int, search, sortBy, order string) ([]models.Product, int, error)
	GetByID(id uint) (*models.Product, error)
	Update(id uint, product *models.Product) error
	Delete(id uint) error
	UpdateStock(id uint, newStock int) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository() ProductRepository {
	return &productRepository{
		db: config.DB,
	}
}

func (r *productRepository) Create(product *models.Product) error {
	query := `
		INSERT INTO products (name, price, stock, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(
		query,
		product.Name,
		product.Price,
		product.Stock,
		now,
		now,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	return err
}

func (r *productRepository) GetAll(page, limit int, search, sortBy, order string) ([]models.Product, int, error) {
	// Default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	offset := (page - 1) * limit

	// Base query
	query := `
		SELECT id, name, price, stock, created_at, updated_at 
		FROM products 
		WHERE deleted_at IS NULL
	`
	countQuery := `
		SELECT COUNT(*) 
		FROM products 
		WHERE deleted_at IS NULL
	`

	// Filtering (search by name)
	var args []interface{}
	if search != "" {
		query += " AND name ILIKE $1"
		countQuery += " AND name ILIKE $1"
		args = append(args, "%"+search+"%")
	}

	// Sorting + pagination
	query += fmt.Sprintf(" ORDER BY %s %s LIMIT %d OFFSET %d", sortBy, order, limit, offset)

	// Ambil data
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	// Hitung total data (tanpa limit/offset)
	var total int
	err = r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}


func (r *productRepository) GetByID(id uint) (*models.Product, error) {
	query := `
		SELECT id, name, price, stock, created_at, updated_at 
		FROM products 
		WHERE id = $1 AND deleted_at IS NULL`

	var product models.Product
	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Update(id uint, product *models.Product) error {
	query := `
		UPDATE products 
		SET name = $1, price = $2, stock = $3, updated_at = $4 
		WHERE id = $5 AND deleted_at IS NULL`

	_, err := r.db.Exec(
		query,
		product.Name,
		product.Price,
		product.Stock,
		time.Now(),
		id,
	)

	return err
}

func (r *productRepository) Delete(id uint) error {
	query := `
		UPDATE products 
		SET deleted_at = $1 
		WHERE id = $2 AND deleted_at IS NULL`

	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

func (r *productRepository) UpdateStock(id uint, newStock int) error {
	query := `
		UPDATE products 
		SET stock = $1, updated_at = $2 
		WHERE id = $3 AND deleted_at IS NULL`

	_, err := r.db.Exec(query, newStock, time.Now(), id)
	return err
}
