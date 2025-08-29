package repositories

import (
	"database/sql"
	"fmt"
	"time"
	"transaction-service/config"
	"transaction-service/models"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetAll(page, limit int, search, sortBy, order string) ([]models.Transaction, int, error)
	GetByID(id uint) (*models.Transaction, error)
	GetTransactionItems(transactionID uint) ([]models.TransactionItem, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{
		db: config.DB,
	}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO transactions (transaction_date, total_amount, created_at, updated_at) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err = tx.QueryRow(
		query,
		transaction.TransactionDate,
		transaction.TotalAmount,
		now,
		now,
	).Scan(&transaction.ID, &transaction.CreatedAt, &transaction.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	for i := range transaction.TransactionItems {
		item := &transaction.TransactionItems[i]
		item.TransactionID = transaction.ID

		itemQuery := `
			INSERT INTO transaction_items (transaction_id, product_id, quantity, subtotal, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING id, created_at, updated_at`

		err = tx.QueryRow(
			itemQuery,
			item.TransactionID,
			item.ProductID,
			item.Quantity,
			item.Subtotal,
			now,
			now,
		).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert transaction item: %w", err)
		}

		updateStockQuery := `
			UPDATE products 
			SET stock = stock - $1, updated_at = $2
			WHERE id = $3 AND stock >= $1 AND deleted_at IS NULL`

		res, err := tx.Exec(updateStockQuery, item.Quantity, now, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to update stock for product_id %d: %w", item.ProductID, err)
		}

		affected, _ := res.RowsAffected()
		if affected == 0 {
			return fmt.Errorf("insufficient stock for product_id %d", item.ProductID)
		}
	}

	return tx.Commit()
}

func (r *transactionRepository) GetAll(page, limit int, search, sortBy, order string) ([]models.Transaction, int, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	if sortBy == "" {
		sortBy = "created_at"
	}
	if order == "" {
		order = "DESC"
	}

	// Validate sort parameters
	allowedSort := map[string]bool{"created_at": true, "transaction_date": true, "total_amount": true, "id": true}
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	// Build search conditions (count uses $1; main uses $3 because main has LIMIT $1 OFFSET $2)
	searchConditionMain := ""
	searchConditionCount := ""
	countArgs := []interface{}{}
	queryArgs := []interface{}{limit, offset}

	if search != "" {
		like := "%" + search + "%"

		// NOTE: saya asumsikan nama kolom produk di tabel products adalah "name".
		// Jika di DB kamu berbeda (mis. "product_name"), ganti p.name -> p.product_name.
		searchConditionCount = `AND (
			CAST(t.id AS TEXT) ILIKE $1
			OR CAST(t.transaction_date AS TEXT) ILIKE $1
			OR CAST(t.total_amount AS TEXT) ILIKE $1
			OR EXISTS (
				SELECT 1 FROM transaction_items ti
				JOIN products p ON p.id = ti.product_id
				WHERE ti.transaction_id = t.id
				  AND p.name ILIKE $1
			)
		)`

		searchConditionMain = `AND (
			CAST(t.id AS TEXT) ILIKE $3
			OR CAST(t.transaction_date AS TEXT) ILIKE $3
			OR CAST(t.total_amount AS TEXT) ILIKE $3
			OR EXISTS (
				SELECT 1 FROM transaction_items ti
				JOIN products p ON p.id = ti.product_id
				WHERE ti.transaction_id = t.id
				  AND p.name ILIKE $3
			)
		)`

		countArgs = append(countArgs, like)
		queryArgs = append(queryArgs, like) // becomes $3 in main query
	}

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(t.id)
		FROM transactions t
		WHERE t.deleted_at IS NULL %s`, searchConditionCount)

	var totalCount int
	if err := r.db.QueryRow(countQuery, countArgs...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT id, transaction_date, total_amount, created_at, updated_at
		FROM transactions t
		WHERE t.deleted_at IS NULL %s
		ORDER BY %s %s
		LIMIT $1 OFFSET $2`, searchConditionMain, sortBy, order)

	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.TransactionDate,
			&transaction.TotalAmount,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Get transaction items (as you had)
		items, err := r.GetTransactionItems(transaction.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get transaction items: %w", err)
		}
		transaction.TransactionItems = items

		transactions = append(transactions, transaction)
	}

	return transactions, totalCount, nil
}



func (r *transactionRepository) GetByID(id uint) (*models.Transaction, error) {
	query := `
		SELECT id, transaction_date, total_amount, created_at, updated_at 
		FROM transactions 
		WHERE id = $1 AND deleted_at IS NULL`

	var transaction models.Transaction
	err := r.db.QueryRow(query, id).Scan(
		&transaction.ID,
		&transaction.TransactionDate,
		&transaction.TotalAmount,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get transaction items
	items, err := r.GetTransactionItems(transaction.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction items: %w", err)
	}
	transaction.TransactionItems = items

	return &transaction, nil
}

func (r *transactionRepository) GetTransactionItems(transactionID uint) ([]models.TransactionItem, error) {
	query := `
		SELECT id, transaction_id, product_id, quantity, subtotal, created_at, updated_at 
		FROM transaction_items 
		WHERE transaction_id = $1 
		ORDER BY created_at ASC`

	rows, err := r.db.Query(query, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TransactionItem
	for rows.Next() {
		var item models.TransactionItem
		err := rows.Scan(
			&item.ID,
			&item.TransactionID,
			&item.ProductID,
			&item.Quantity,
			&item.Subtotal,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}