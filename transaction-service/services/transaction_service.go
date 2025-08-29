package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
	"transaction-service/clients"
	"transaction-service/dto"
	"transaction-service/models"
	"transaction-service/repositories"
)

type TransactionService interface {
	CreateTransaction(req *dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	GetAllTransactions(page, limit int, search, sortBy, order string) ([]dto.TransactionResponse, int, error)
	GetTransactionByID(id uint) (*dto.TransactionResponse, error)
}

type transactionService struct {
	repo          repositories.TransactionRepository
	productClient clients.ProductClient
}

func NewTransactionService(repo repositories.TransactionRepository, productClient clients.ProductClient) TransactionService {
	return &transactionService{
		repo:          repo,
		productClient: productClient,
	}
}

func (s *transactionService) CreateTransaction(req *dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("transaction must have at least one item")
	}

	// Validate request items
	for _, item := range req.Items {
		if item.ProductID == 0 {
			return nil, errors.New("product ID is required")
		}
		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than 0")
		}
	}

	transaction := &models.Transaction{
		TransactionDate: time.Now(),
		TotalAmount:     0,
	}

	// Process each item and calculate total
	var totalAmount float64
	for _, item := range req.Items {
		// Get product details from product service
		product, err := s.productClient.GetByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product with ID %d not found or service unavailable", item.ProductID)
		}

		// Check stock availability
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product '%s'. Available: %d, Requested: %d", 
				product.Name, product.Stock, item.Quantity)
		}

		// Calculate subtotal
		subtotal := product.Price * float64(item.Quantity)

		// Create transaction item
		transactionItem := models.TransactionItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
		}

		transaction.TransactionItems = append(transaction.TransactionItems, transactionItem)
		totalAmount += subtotal
	}

	transaction.TotalAmount = totalAmount

	// Save transaction
	err := s.repo.Create(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return s.modelToResponse(transaction)
}

func (s *transactionService) GetAllTransactions(page, limit int, search, sortBy, order string) ([]dto.TransactionResponse, int, error) {
	transactions, total, err := s.repo.GetAll(page, limit, search, sortBy, order)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.TransactionResponse
	for _, transaction := range transactions {
		response, err := s.modelToResponseWithFallback(&transaction)
		if err != nil {
			log.Printf("Warning: Failed to get complete product details for transaction %d: %v", transaction.ID, err)
			// Continue with partial data instead of failing completely
			continue
		}
		responses = append(responses, *response)
	}

	return responses, total, nil
}

func (s *transactionService) GetTransactionByID(id uint) (*dto.TransactionResponse, error) {
	transaction, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	return s.modelToResponseWithFallback(transaction)
}

func (s *transactionService) modelToResponse(transaction *models.Transaction) (*dto.TransactionResponse, error) {
	var items []dto.TransactionItemResponse

	// Collect unique product IDs
	productIDs := make([]uint, 0, len(transaction.TransactionItems))
	for _, item := range transaction.TransactionItems {
		productIDs = append(productIDs, item.ProductID)
	}

	// Fetch product details in batch
	products, err := s.productClient.GetMultiple(productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product details: %w", err)
	}

	// Build response items with product details
	for _, item := range transaction.TransactionItems {
		product, exists := products[item.ProductID]
		if !exists {
			return nil, fmt.Errorf("product details not found for product ID %d", item.ProductID)
		}

		items = append(items, dto.TransactionItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    item.Quantity,
			Subtotal:    item.Subtotal,
		})
	}

	return &dto.TransactionResponse{
		ID:               transaction.ID,
		TransactionDate:  transaction.TransactionDate.Format("2006-01-02 15:04:05"),
		TotalAmount:      transaction.TotalAmount,
		TransactionItems: items,
		CreatedAt:        transaction.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// New method with fallback for deleted products
func (s *transactionService) modelToResponseWithFallback(transaction *models.Transaction) (*dto.TransactionResponse, error) {
	var items []dto.TransactionItemResponse

	// Collect unique product IDs
	productIDs := make([]uint, 0, len(transaction.TransactionItems))
	for _, item := range transaction.TransactionItems {
		productIDs = append(productIDs, item.ProductID)
	}

	// Try to fetch product details
	products := make(map[uint]*clients.ProductResponse)
	
	// Fetch products individually to handle deleted products gracefully
	for _, productID := range productIDs {
		product, err := s.productClient.GetByID(productID)
		if err != nil {
			// Log warning but continue with fallback data
			log.Printf("Warning: Could not fetch product %d: %v", productID, err)
		} else {
			products[productID] = product
		}
	}

	// Build response items with fallback for deleted products
	for _, item := range transaction.TransactionItems {
		var productName string
		var productPrice float64

		if product, exists := products[item.ProductID]; exists {
			// Product still exists, use current data
			productName = product.Name
			productPrice = product.Price
		} else {
			// Product deleted or unavailable, use fallback
			productName = fmt.Sprintf("Product ID %d (Deleted)", item.ProductID)
			productPrice = item.Subtotal / float64(item.Quantity) // Calculate from subtotal
		}

		items = append(items, dto.TransactionItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: productName,
			Price:       productPrice,
			Quantity:    item.Quantity,
			Subtotal:    item.Subtotal,
		})
	}

	return &dto.TransactionResponse{
		ID:               transaction.ID,
		TransactionDate:  transaction.TransactionDate.Format("2006-01-02 15:04:05"),
		TotalAmount:      transaction.TotalAmount,
		TransactionItems: items,
		CreatedAt:        transaction.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}