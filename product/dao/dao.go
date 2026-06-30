package dao

import "github.com/saiteja/ecommerce/product/models"

type ProductDAO interface {
	CreateProduct(product *models.Product) (*models.Product, error)
	GetProductByID(id string) (*models.Product, error)
	GetProducts(page, pageSize int) ([]*models.Product, int, error)
}

type InventoryDAO interface {
	GetInventory(productID string) (*models.ProductInventory, error)
	SetInventory(inventory *models.ProductInventory) error
}
