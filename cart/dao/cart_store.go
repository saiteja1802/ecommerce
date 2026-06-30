package dao

import (
	"errors"
	"sync"

	"github.com/mohae/deepcopy"
	"github.com/saiteja/ecommerce/cart/models"
)

var (
	ErrCartNotFound = errors.New("cart not found")
	ErrItemNotFound = errors.New("item not found in cart")
)

type InMemoryCartStore struct {
	mu    sync.RWMutex
	carts map[string]*models.Cart // userID -> Cart
}

func NewInMemoryCartStore() *InMemoryCartStore {
	return &InMemoryCartStore{carts: make(map[string]*models.Cart)}
}

func (s *InMemoryCartStore) GetCart(userID string) (*models.Cart, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cart, ok := s.carts[userID]
	if !ok {
		return nil, ErrCartNotFound
	}
	// Return a deep copy so the caller can iterate Items after the lock is
	// released without racing against concurrent SetItem / RemoveItem calls.
	return deepcopy.Copy(cart).(*models.Cart), nil
}

func (s *InMemoryCartStore) SetItem(userID, productID string, quantity int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, ok := s.carts[userID]
	if !ok {
		cart = &models.Cart{UserID: userID}
		s.carts[userID] = cart
	}
	for _, item := range cart.GetItems() {
		if item.GetProductID() == productID {
			item.Quantity = quantity
			return nil
		}
	}
	cart.Items = append(cart.GetItems(), &models.CartItem{
		ProductID: productID,
		Quantity:  quantity,
	})
	return nil
}

func (s *InMemoryCartStore) RemoveItem(userID, productID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, ok := s.carts[userID]
	if !ok {
		return ErrItemNotFound
	}
	for i, item := range cart.GetItems() {
		if item.GetProductID() == productID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			return nil
		}
	}
	return ErrItemNotFound
}
