package services

import (
	"database/sql"
	"errors"
	"product-service/dto"
	"product-service/models"
	"product-service/repositories"
)

type ProductService interface {
	CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetAllProducts() ([]dto.ProductResponse, error)
	GetProductByID(id uint) (*dto.ProductResponse, error)
	UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(id uint) error
	UpdateStock(id uint, newStock int) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := &models.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	err := s.repo.Create(product)
	if err != nil {
		return nil, err
	}

	return s.modelToResponse(product), nil
}

func (s *productService) GetAllProducts() ([]dto.ProductResponse, error) {
	products, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.ProductResponse
	for _, product := range products {
		responses = append(responses, *s.modelToResponse(&product))
	}

	return responses, nil
}

func (s *productService) GetProductByID(id uint) (*dto.ProductResponse, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return s.modelToResponse(product), nil
}

func (s *productService) UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	// Get existing product
	existingProduct, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Update kolom yang diubah saja
	updateData := &models.Product{
		Name:  existingProduct.Name,
		Price: existingProduct.Price,
		Stock: existingProduct.Stock,
	}

	if req.Name != "" {
		updateData.Name = req.Name
	}
	if req.Price > 0 {
		updateData.Price = req.Price
	}
	if req.Stock >= 0 {
		updateData.Stock = req.Stock
	}

	err = s.repo.Update(id, updateData)
	if err != nil {
		return nil, err
	}

	// ambil produk yang sudah diupdate untuk response
	updatedProduct, _ := s.repo.GetByID(id)
	return s.modelToResponse(updatedProduct), nil
}

func (s *productService) DeleteProduct(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("product not found")
		}
		return err
	}

	return s.repo.Delete(id)
}

func (s *productService) UpdateStock(id uint, newStock int) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("product not found")
		}
		return err
	}

	return s.repo.UpdateStock(id, newStock)
}

func (s *productService) modelToResponse(product *models.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
