package main

import (
	"net/http"

	"github.com/govalues/decimal"
	"github.com/govalues/money"
	"github.com/saiteja/ecommerce/auth"
	authdao "github.com/saiteja/ecommerce/auth/dao"
	"github.com/saiteja/ecommerce/cart"
	cartdao "github.com/saiteja/ecommerce/cart/dao"
	cartmodels "github.com/saiteja/ecommerce/cart/models"
	"github.com/saiteja/ecommerce/cmd/httpserver"
	"github.com/saiteja/ecommerce/pkg/logger"
	"github.com/saiteja/ecommerce/product"
	productdao "github.com/saiteja/ecommerce/product/dao"
	productmodels "github.com/saiteja/ecommerce/product/models"
)

func main() {
	credentialsStore := authdao.NewInMemoryUserCredentialsStore()
	sessionStore := authdao.NewInMemorySessionStore()
	authService := auth.NewService(credentialsStore, sessionStore)

	productStore := productdao.NewInMemoryProductStore()
	inventoryStore := productdao.NewInMemoryInventoryStore()
	productService := product.NewService(productStore, inventoryStore)

	cartStore := cartdao.NewInMemoryCartStore()
	couponStore := cartdao.NewInMemoryCouponStore()
	cartService := cart.NewService(cartStore, couponStore, productService)

	s := httpserver.New(authService, productService, cartService, productStore, inventoryStore, couponStore)
	seedData(s)

	logger.L.Info("server started", "addr", ":8080")
	if err := http.ListenAndServe(":8080", s); err != nil {
		logger.L.Error("server stopped", "error", err)
	}
}

func seedData(s *httpserver.Server) {
	catalogue := []struct {
		name, description string
		paise             int64
		stock             int
	}{
		{"Laptop", "High-performance laptop", 79999, 20},
		{"Phone", "Latest smartphone", 49999, 50},
		{"Headphones", "Noise-cancelling headphones", 19999, 100},
	}
	for _, item := range catalogue {
		p, err := s.ProductStore.CreateProduct(&productmodels.Product{
			Name:        item.name,
			Description: item.description,
			Price:       money.MustNewAmount("INR", item.paise, cartmodels.InrScale),
		})
		if err != nil {
			logger.L.Error("seed: failed to create product", "name", item.name, "error", err)
			continue
		}
		if err := s.InventoryStore.SetInventory(&productmodels.ProductInventory{
			ProductID: p.GetID(),
			Quantity:  item.stock,
		}); err != nil {
			logger.L.Error("seed: failed to set inventory", "name", item.name, "error", err)
		}
	}

	coupons := []struct {
		name       string
		percentage string
		maxPaise   int64
	}{
		{"SAVE10", "0.10", 50000},  // 10% off, max 500.00
		{"WELCOME20", "0.20", 20000}, // 20% off, max 200.00
	}
	for _, c := range coupons {
		maxDiscount, _ := money.NewAmount("INR", c.maxPaise, cartmodels.InrScale)
		if err := s.CouponStore.CreateCoupon(&cartmodels.Coupon{
			Name:               c.name,
			DiscountPercentage: decimal.MustParse(c.percentage),
			MaxDiscount:        maxDiscount,
		}); err != nil {
			logger.L.Error("seed: failed to create coupon", "name", c.name, "error", err)
		}
	}

	logger.L.Info("sample data loaded", "products", len(catalogue), "coupons", len(coupons))
}
