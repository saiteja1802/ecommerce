package dao

import (
	"errors"
	"sync"

	"github.com/saiteja/ecommerce/cart/models"
)

var ErrCouponNotFound = errors.New("coupon not found")

type CouponDAO interface {
	GetCoupon(name string) (*models.Coupon, error)
	CreateCoupon(coupon *models.Coupon) error
}

type InMemoryCouponStore struct {
	mu      sync.RWMutex
	coupons map[string]*models.Coupon // name -> Coupon
}

func NewInMemoryCouponStore() *InMemoryCouponStore {
	return &InMemoryCouponStore{coupons: make(map[string]*models.Coupon)}
}

func (s *InMemoryCouponStore) GetCoupon(name string) (*models.Coupon, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	coupon, ok := s.coupons[name]
	if !ok {
		return nil, ErrCouponNotFound
	}
	return coupon, nil
}

func (s *InMemoryCouponStore) CreateCoupon(coupon *models.Coupon) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.coupons[coupon.Name] = coupon
	return nil
}
