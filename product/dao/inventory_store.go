package dao

import (
	"sync"

	"github.com/saiteja/ecommerce/product/models"
)

type InMemoryInventoryStore struct {
	mu    sync.RWMutex
	stock map[string]int // productID -> quantity
}

func NewInMemoryInventoryStore() *InMemoryInventoryStore {
	return &InMemoryInventoryStore{stock: make(map[string]int)}
}

func (s *InMemoryInventoryStore) GetInventory(productID string) (*models.ProductInventory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	qty, ok := s.stock[productID]
	if !ok {
		return nil, ErrProductNotFound
	}
	return &models.ProductInventory{
		ProductID: productID,
		Quantity:  qty,
	}, nil
}

func (s *InMemoryInventoryStore) SetInventory(inventory *models.ProductInventory) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stock[inventory.GetProductID()] = inventory.GetQuantity()
	return nil
}
