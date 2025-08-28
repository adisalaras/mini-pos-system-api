package repositories

import (
	"database/sql"
	"time"
	"transaction-service/config"
	"transaction-service/models"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetAll() ([]models.Transaction, error)
	GetByID(id uint) (*models.Transaction, error)
	CreateTransactionItem(item *models.TransactionItem) error
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
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert transaction
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
		return err
	}

	// Insert transaction items
	for i := range transaction.TransactionItems {
		item := &transaction.TransactionItems[i]
		item.TransactionID = transaction.ID

		itemQuery := `
			INSERT INTO transaction_items (transaction_id, product_id, product_name, price, quantity, subtotal, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
			RETURNING id, created_at, updated_at`

		err = tx.QueryRow(
			itemQuery,
			item.TransactionID,
			item.ProductID,
			item.ProductName,
			item.Price,
			item.Quantity,
			item.Subtotal,
			now,
			now,
		).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *transactionRepository) GetTransactionItems(transactionID uint) ([]models.TransactionItem, error) {
	query := `
		SELECT id, transaction_id, product_id, product_name, price, quantity, subtotal, created_at, updated_at 
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
			&item.ProductName,
			&item.Price,
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

func (r *transactionRepository) GetAll() ([]models.Transaction, error) {
	query := `
	SELECT id, transaction_date, total_amount, created_at, updated_at
	FROM transactions
	WHERE deleted_at IS NULL
	ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		items, err := r.GetTransactionItems(transaction.ID)
		if err != nil {
			return nil, err
		}
		transaction.TransactionItems = items
		transactions = append(transactions, transaction)
	}
	return transactions, nil
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

	items, err := r.GetTransactionItems(transaction.ID)
	if err != nil {
		return nil, err
	}
	transaction.TransactionItems = items

	return &transaction, nil
}

func (r *transactionRepository) CreateTransactionItem(item *models.TransactionItem) error {
	query := `
		INSERT INTO transaction_items (transaction_id, product_id, product_name, price, quantity, subtotal, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(
		query,
		item.TransactionID,
		item.ProductID,
		item.ProductName,
		item.Price,
		item.Quantity,
		item.Subtotal,
		now,
		now,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	return err
}
