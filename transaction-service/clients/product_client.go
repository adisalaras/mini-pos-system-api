package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ProductResponse struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

type ApiResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    ProductResponse `json:"data"`
}

type ProductClient interface {
	GetByID(id uint) (*ProductResponse, error)
	GetMultiple(ids []uint) (map[uint]*ProductResponse, error)
	GetByIDWithFallback(id uint) (*ProductResponse, bool) // Returns product and exists flag
}

type productClient struct {
	baseURL string
	client  *http.Client
}

func NewProductClient(baseURL string) ProductClient {
	return &productClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *productClient) GetByID(id uint) (*ProductResponse, error) {
	url := fmt.Sprintf("%s/api/products/%d", c.baseURL, id)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call product service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product with ID %d not found", id)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product service returned status %d", resp.StatusCode)
	}

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("product service error: %s", apiResp.Message)
	}

	return &apiResp.Data, nil
}

func (c *productClient) GetByIDWithFallback(id uint) (*ProductResponse, bool) {
	product, err := c.GetByID(id)
	if err != nil {
		log.Printf("Product %d not found: %v", id, err)
		return nil, false
	}
	return product, true
}

func (c *productClient) GetMultiple(ids []uint) (map[uint]*ProductResponse, error) {
	products := make(map[uint]*ProductResponse)
	failedIDs := make([]uint, 0)

	// Fetch each product individually with error tolerance
	for _, id := range ids {
		product, exists := c.GetByIDWithFallback(id)
		if exists {
			products[id] = product
		} else {
			failedIDs = append(failedIDs, id)
		}
	}

	// Log warning for failed products but don't return error
	if len(failedIDs) > 0 {
		log.Printf("Warning: Could not fetch products with IDs: %v", failedIDs)
	}

	return products, nil
}
