package dao

import (
	"errors"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/saiteja/ecommerce/product/models"
)

var ErrProductNotFound = errors.New("product not found")

type InMemoryProductStore struct {
	mu       sync.RWMutex
	products []*models.Product
}

func NewInMemoryProductStore() *InMemoryProductStore {
	return &InMemoryProductStore{}
}

func (s *InMemoryProductStore) CreateProduct(product *models.Product) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	product.ID = "PRD" + ulid.Make().String()
	if product.GetCreatedAt().IsZero() {
		product.CreatedAt = time.Now().UTC()
	}

	// Copy so the store's record is independent of the caller's pointer.
	cp := *product
	s.products = append(s.products, &cp)
	return product, nil
}

func (s *InMemoryProductStore) GetProductByID(id string) (*models.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, p := range s.products {
		if p.GetID() == id {
			// Copy so the caller cannot mutate the store's internal record.
			cp := *p
			return &cp, nil
		}
	}
	return nil, ErrProductNotFound
}

func (s *InMemoryProductStore) GetProducts(page, pageSize int) ([]*models.Product, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := len(s.products)
	start := (page - 1) * pageSize
	if start >= total {
		return []*models.Product{}, total, nil
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	result := make([]*models.Product, end-start)
	for i, p := range s.products[start:end] {
		// Copy so the caller cannot mutate the store's internal record.
		cp := *p
		result[i] = &cp
	}
	return result, total, nil
}
