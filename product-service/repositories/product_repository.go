package repositories

import (
	"database/sql"
	"product-service/config"
	"product-service/models"
	"time"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetAll() ([]models.Product, error)
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

func (r *productRepository) GetAll() ([]models.Product, error) {
	query := `
		SELECT id, name, price, stock, created_at, updated_at 
		FROM products 
		WHERE deleted_at IS NULL 
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
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
