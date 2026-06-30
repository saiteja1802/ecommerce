package integration_test

import (
	"net/http"
	"testing"

	"github.com/govalues/money"
	productapi "github.com/saiteja/ecommerce/product/api_models"
	productmodels "github.com/saiteja/ecommerce/product/models"
)

func TestGetProductCatalog(t *testing.T) {
	ts, cleanup := startServer(t)
	defer cleanup()

	ts.ProductStore.CreateProduct(&productmodels.Product{Name: "Laptop", Description: "Gaming laptop", Price: money.MustNewAmount("INR", 99999, 2)})
	ts.ProductStore.CreateProduct(&productmodels.Product{Name: "Phone", Description: "Smartphone", Price: money.MustNewAmount("INR", 49999, 2)})
	ts.ProductStore.CreateProduct(&productmodels.Product{Name: "Headphones", Description: "Noise cancelling", Price: money.MustNewAmount("INR", 19999, 2)})

	token := signupAndLogin(t, ts, "buyer@example.com", "correct-password")

	catalog := &productapi.GetProductsCatalogResponse{}
	get(t, ts, "/products?page=1&page_size=2", token, http.StatusOK, catalog)

	if catalog.GetTotalProducts() != 3 {
		t.Fatalf("expected total 3, got %d", catalog.GetTotalProducts())
	}
	if len(catalog.GetProducts()) != 2 {
		t.Fatalf("expected 2 products on page 1, got %d", len(catalog.GetProducts()))
	}
	if catalog.GetPage() != 1 {
		t.Fatalf("expected page 1, got %d", catalog.GetPage())
	}
	if catalog.GetPageSize() != 2 {
		t.Fatalf("expected page_size 2, got %d", catalog.GetPageSize())
	}
}

func TestGetProductDetails(t *testing.T) {
	ts, cleanup := startServer(t)
	defer cleanup()

	p, _ := ts.ProductStore.CreateProduct(&productmodels.Product{
		Name:        "Laptop",
		Description: "Gaming laptop",
		Price:       money.MustNewAmount("INR", 99999, 2),
	})

	token := signupAndLogin(t, ts, "buyer@example.com", "correct-password")

	resp := &productapi.GetProductDetailsResponse{}
	get(t, ts, "/products/"+p.GetID(), token, http.StatusOK, resp)

	if resp.GetID() != p.GetID() {
		t.Fatalf("expected id %q, got %q", p.GetID(), resp.GetID())
	}
	if resp.GetName() != "Laptop" {
		t.Fatalf("expected name %q, got %q", "Laptop", resp.GetName())
	}
	if resp.GetDescription() != "Gaming laptop" {
		t.Fatalf("expected description %q, got %q", "Gaming laptop", resp.GetDescription())
	}
	if resp.GetPrice() != "999.99" {
		t.Fatalf("expected price %q, got %q", "999.99", resp.GetPrice())
	}
	if resp.GetCurrencyCode() != "INR" {
		t.Fatalf("expected currency_code %q, got %q", "INR", resp.GetCurrencyCode())
	}
}

func TestGetProductDetailsNotFound(t *testing.T) {
	ts, cleanup := startServer(t)
	defer cleanup()

	token := signupAndLogin(t, ts, "buyer@example.com", "correct-password")
	get(t, ts, "/products/nonexistent-id", token, http.StatusNotFound, nil)
}

func TestGetProductInventory(t *testing.T) {
	ts, cleanup := startServer(t)
	defer cleanup()

	p, _ := ts.ProductStore.CreateProduct(&productmodels.Product{
		Name:        "Laptop",
		Description: "Gaming laptop",
		Price:       money.MustNewAmount("INR", 99999, 2),
	})
	ts.InventoryStore.SetInventory(&productmodels.ProductInventory{ProductID: p.GetID(), Quantity: 42})

	token := signupAndLogin(t, ts, "buyer@example.com", "correct-password")

	resp := &productapi.GetInventoryResponse{}
	get(t, ts, "/products/"+p.GetID()+"/inventory", token, http.StatusOK, resp)

	if resp.GetProductID() != p.GetID() {
		t.Fatalf("expected product_id %q, got %q", p.GetID(), resp.GetProductID())
	}
	if resp.GetQuantity() != 42 {
		t.Fatalf("expected quantity 42, got %d", resp.GetQuantity())
	}
}
