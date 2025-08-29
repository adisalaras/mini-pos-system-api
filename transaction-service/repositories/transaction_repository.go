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

	// Insert transaction items & update stock
	for i := range transaction.TransactionItems {
		item := &transaction.TransactionItems[i]
		item.TransactionID = transaction.ID

		// 1. Insert transaction item
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

		// 2. Update stock produk (efisien: tanpa RETURNING)
		updateStockQuery := `
			UPDATE products 
			SET stock = stock - $1, updated_at = $2
			WHERE id = $3 AND stock >= $1`

		res, err := tx.Exec(updateStockQuery, item.Quantity, now, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to update stock for product_id %d: %w", item.ProductID, err)
		}

		// Cek apakah ada row yang ter-update (stok cukup)
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return fmt.Errorf("insufficient stock for product_id %d", item.ProductID)
		}
	}

	// Commit transaction
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

func (r *transactionRepository) GetAll(page, limit int, search, sortBy, order string) ([]models.Transaction, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Default sorting
	if sortBy == "" {
		sortBy = "created_at"
	}
	if order == "" {
		order = "DESC"
	}

	// Allowed sort columns
	allowedSort := map[string]bool{"created_at": true, "transaction_date": true, "total_amount": true}
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	query := fmt.Sprintf(`
		SELECT  t.id, t.transaction_date, t.total_amount, t.created_at, t.updated_at, 
		ti.id, ti.product_id, ti.product_name, ti.price, ti.quantity, ti.subtotal,
			COUNT(*) OVER() as total_count
		FROM transactions t
		LEFT JOIN transaction_items ti ON ti.transaction_id = t.id
		WHERE t.deleted_at IS NULL
		%s
		ORDER BY t.%s %s
		LIMIT $1 OFFSET $2
	`,
		// Filtering (search)
		func() string {
			if search != "" {
				return "AND (ti.product_name ILIKE '%' || $3 || '%')"
			}
			return ""
		}(),
		sortBy, order,
	)

	var rows *sql.Rows
	var err error
	if search != "" {
		rows, err = r.db.Query(query, limit, offset, search)
	} else {
		rows, err = r.db.Query(query, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	txMap := make(map[uint]*models.Transaction)

	totalCount := 0

	for rows.Next() {
		var (
			tx              models.Transaction
			item            models.TransactionItem
			nullItemID      sql.NullInt64
			nullProductID   sql.NullInt64
			nullProductName sql.NullString
			nullPrice       sql.NullFloat64
			nullQty         sql.NullInt64
			nullSubtotal    sql.NullFloat64
		)

		err := rows.Scan(
			&tx.ID, &tx.TransactionDate, &tx.TotalAmount, &tx.CreatedAt, &tx.UpdatedAt,
			&nullItemID, &nullProductID, &nullProductName, &nullPrice, &nullQty, &nullSubtotal,
			&totalCount,
		)
		if err != nil {
			return nil, 0, err
		}

		// Masukkan ke map biar tidak duplikat transaksi
		if _, exists := txMap[tx.ID]; !exists {
			txMap[tx.ID] = &tx
		}

		// Tambahkan item kalau ada
		if nullItemID.Valid {
			item.ID = uint(nullItemID.Int64)
			item.TransactionID = tx.ID
			item.ProductID = uint(nullProductID.Int64)
			item.ProductName = nullProductName.String
			item.Price = nullPrice.Float64
			item.Quantity = int(nullQty.Int64)
			item.Subtotal = nullSubtotal.Float64
			txMap[tx.ID].TransactionItems = append(txMap[tx.ID].TransactionItems, item)
		}
	}

	transactions := make([]models.Transaction, 0, len(txMap))
	for _, t := range txMap {
		transactions = append(transactions, *t)
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
