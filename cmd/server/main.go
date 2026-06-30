package main

import (
	"net/http"

	"github.com/saiteja/ecommerce/auth"
	authdao "github.com/saiteja/ecommerce/auth/dao"
	"github.com/saiteja/ecommerce/cmd/httpserver"
	"github.com/saiteja/ecommerce/pkg/logger"
	"github.com/saiteja/ecommerce/product"
	productdao "github.com/saiteja/ecommerce/product/dao"
)

func main() {
	credentialsStore := authdao.NewInMemoryUserCredentialsStore()
	sessionStore := authdao.NewInMemorySessionStore()
	authService := auth.NewService(credentialsStore, sessionStore)

	productStore := productdao.NewInMemoryProductStore()
	inventoryStore := productdao.NewInMemoryInventoryStore()
	productService := product.NewService(productStore, inventoryStore)

	s := httpserver.New(authService, productService, productStore, inventoryStore)

	logger.L.Info("server started", "addr", ":8080")
	if err := http.ListenAndServe(":8080", s); err != nil {
		logger.L.Error("server stopped", "error", err)
	}
}
