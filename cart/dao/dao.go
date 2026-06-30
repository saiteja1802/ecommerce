package dao

import "github.com/saiteja/ecommerce/cart/models"

type CartDAO interface {
	GetCart(userID string) (*models.Cart, error)
	SetItem(userID, productID string, quantity int) error
	RemoveItem(userID, productID string) error
}
