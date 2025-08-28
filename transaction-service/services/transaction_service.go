package services

import (
	"database/sql"
	"errors"
	"time"
	"transaction-service/dto"
	"transaction-service/models"
	"transaction-service/repositories"
)

type TransactionService interface {
	CreateTransaction(req *dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	GetAllTransactions() ([]dto.TransactionResponse, error)
	GetTransactionByID(id uint) (*dto.TransactionResponse, error)
}

type transactionService struct {
	repo repositories.TransactionRepository
}

func NewTransactionService(repo repositories.TransactionRepository) TransactionService {
	return &transactionService{
		repo: repo,
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
		if item.ProductName == "" {
			return nil, errors.New("product name is required")
		}
		if item.Price <= 0 {
			return nil, errors.New("price must be greater than 0")
		}
		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than 0")
		}
	}

	transaction := &models.Transaction{
		TransactionDate: time.Now(),
		TotalAmount:     0,
	}

	// Create transaction items and calculate total
	var totalAmount float64
	for _, item := range req.Items {
		subtotal := item.Price * float64(item.Quantity)

		transactionItem := models.TransactionItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		}

		transaction.TransactionItems = append(transaction.TransactionItems, transactionItem)
		totalAmount += subtotal
	}

	transaction.TotalAmount = totalAmount

	err := s.repo.Create(transaction)
	if err != nil {
		return nil, err
	}

	return s.modelToResponse(transaction), nil
}

func (s *transactionService) GetAllTransactions() ([]dto.TransactionResponse, error) {
	transactions, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.TransactionResponse
	for _, transaction := range transactions {
		responses = append(responses, *s.modelToResponse(&transaction))
	}

	return responses, nil
}

func (s *transactionService) GetTransactionByID(id uint) (*dto.TransactionResponse, error) {
	transaction, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	return s.modelToResponse(transaction), nil
}

func (s *transactionService) modelToResponse(transaction *models.Transaction) *dto.TransactionResponse {
	var items []dto.TransactionItemResponse
	for _, item := range transaction.TransactionItems {
		items = append(items, dto.TransactionItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Price:       item.Price,
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
	}
}
