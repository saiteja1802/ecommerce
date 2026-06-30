package integration_test

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/govalues/money"
	cartapi "github.com/saiteja/ecommerce/cart/api_models"
	productmodels "github.com/saiteja/ecommerce/product/models"
)

// seedProduct creates a product and sets its stock, returning the product ID.
// price is in paise (e.g. 999 = INR 9.99).
func seedProduct(t *testing.T, ts *testServer, name string, price int64, stock int) string {
	t.Helper()
	p, err := ts.ProductStore.CreateProduct(&productmodels.Product{Name: name, Price: money.MustNewAmount("INR", price, 2)})
	if err != nil {
		t.Fatalf("create product: %v", err)
	}
	if err := ts.InventoryStore.SetInventory(&productmodels.ProductInventory{ProductID: p.GetID(), Quantity: stock}); err != nil {
		t.Fatalf("set inventory: %v", err)
	}
	return p.GetID()
}

func TestAddItem(t *testing.T) {
	ts, cleanup := startServer(t)
	defer cleanup()

	// adding an item to an empty cart returns the item with the requested quantity
	{
		productID := seedProduct(t, ts, "Laptop", 99999, 10)
		token := signupAndLogin(t, ts, "buyer-add@example.com", "password")
		resp := &cartapi.AddItemResponse{}
		post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: productID, Quantity: 2}, token, http.StatusOK, resp)
		if len(resp.GetItems()) != 1 {
			t.Fatalf("expected 1 item, got %d", len(resp.GetItems()))
		}
		if resp.GetItems()[0].GetQuantity() != 2 {
			t.Fatalf("expected qty 2, got %d", resp.GetItems()[0].GetQuantity())
		}
	}

	// adding the same product twice merges the quantities into a single cart entry
	{
		productID := seedProduct(t, ts, "Laptop", 99999, 10)
		token := signupAndLogin(t, ts, "buyer-merge@example.com", "password")
		post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: productID, Quantity: 2}, token, http.StatusOK, nil)
		resp := &cartapi.AddItemResponse{}
		post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: productID, Quantity: 3}, token, http.StatusOK, resp)
		if resp.GetItems()[0].GetQuantity() != 5 {
			t.Fatalf("expected merged qty 5, got %d", resp.GetItems()[0].GetQuantity())
		}
	}

	// adding more quantity than available stock returns 422 Unprocessable Entity
	{
		productID := seedProduct(t, ts, "Laptop", 99999, 1)
		token := signupAndLogin(t, ts, "buyer-stock@example.com", "password")
		post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: productID, Quantity: 5}, token, http.StatusUnprocessableEntity, nil)
	}

	// quantity of zero is rejected before any product lookup
	{
		token := signupAndLogin(t, ts, "buyer-zerqty@example.com", "password")
		post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: "any-id", Quantity: 0}, token, http.StatusBadRequest, nil)
	}

	// adding a product that does not exist returns 404
	{
		token := signupAndLogin(t, ts, "buyer-notfound@example.com", "password")
		post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: "nonexistent", Quantity: 1}, token, http.StatusNotFound, nil)
	}

	// request without a bearer token is rejected with 401
	{
		post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: "any", Quantity: 1}, "", http.StatusUnauthorized, nil)
	}
}

func TestRemoveItem(t *testing.T) {
	ts, cleanup := startServer(t)
	defer cleanup()

	// removing an item that is in the cart returns the updated (empty) cart
	{
		productID := seedProduct(t, ts, "Laptop", 99999, 10)
		token := signupAndLogin(t, ts, "buyer-remove@example.com", "password")
		post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: productID, Quantity: 2}, token, http.StatusOK, nil)
		resp := &cartapi.RemoveItemResponse{}
		del(t, ts, "/cart/items/"+productID, token, http.StatusOK, resp)
		if len(resp.GetItems()) != 0 {
			t.Fatalf("expected empty cart after removal, got %d items", len(resp.GetItems()))
		}
	}

	// removing a product that was never added to the cart returns 404
	{
		token := signupAndLogin(t, ts, "buyer-removenotfound@example.com", "password")
		del(t, ts, "/cart/items/nonexistent-product", token, http.StatusNotFound, nil)
	}
}

func TestUpdateQuantity(t *testing.T) {
	ts, cleanup := startServer(t)
	defer cleanup()

	// updating quantity of an existing cart item reflects the new quantity in the response
	{
		productID := seedProduct(t, ts, "Laptop", 99999, 10)
		token := signupAndLogin(t, ts, "buyer-update@example.com", "password")
		post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: productID, Quantity: 2}, token, http.StatusOK, nil)
		resp := &cartapi.UpdateQuantityResponse{}
		patch(t, ts, "/cart/items/"+productID, cartapi.UpdateQuantityRequest{Quantity: 7}, token, http.StatusOK, resp)
		if resp.GetItems()[0].GetQuantity() != 7 {
			t.Fatalf("expected qty 7, got %d", resp.GetItems()[0].GetQuantity())
		}
	}

	// updating quantity to 0 removes the item from the cart entirely
	{
		productID := seedProduct(t, ts, "Laptop", 99999, 10)
		token := signupAndLogin(t, ts, "buyer-updatezero@example.com", "password")
		post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: productID, Quantity: 2}, token, http.StatusOK, nil)
		resp := &cartapi.UpdateQuantityResponse{}
		patch(t, ts, "/cart/items/"+productID, cartapi.UpdateQuantityRequest{Quantity: 0}, token, http.StatusOK, resp)
		if len(resp.GetItems()) != 0 {
			t.Fatalf("expected empty cart after qty 0 update, got %d items", len(resp.GetItems()))
		}
	}

	// updating to a quantity greater than available stock returns 422
	{
		productID := seedProduct(t, ts, "Laptop", 99999, 3)
		token := signupAndLogin(t, ts, "buyer-updatestock@example.com", "password")
		post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: productID, Quantity: 1}, token, http.StatusOK, nil)
		patch(t, ts, "/cart/items/"+productID, cartapi.UpdateQuantityRequest{Quantity: 10}, token, http.StatusUnprocessableEntity, nil)
	}
}

func TestCartIsolatedBetweenUsers(t *testing.T) {
	ts, cleanup := startServer(t)
	defer cleanup()
	productID := seedProduct(t, ts, "Laptop", 99999, 10)

	tokenA := signupAndLogin(t, ts, "user-a@example.com", "password-a")
	post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: productID, Quantity: 3}, tokenA, http.StatusOK, nil)

	respA := &cartapi.CartTotalResponse{}
	get(t, ts, "/cart", tokenA, http.StatusOK, respA)
	if len(respA.GetItems()) != 1 || respA.GetItems()[0].GetQuantity() != 3 {
		t.Fatalf("user A should have 1 item with qty 3, got %v", respA.GetItems())
	}

	tokenB := signupAndLogin(t, ts, "user-b@example.com", "password-b")
	respB := &cartapi.CartTotalResponse{}
	get(t, ts, "/cart", tokenB, http.StatusOK, respB)
	if len(respB.GetItems()) != 0 {
		t.Fatalf("user B should see an empty cart, got %d items", len(respB.GetItems()))
	}
}

// TestGetCartTotalMoneyMath verifies that totals are computed without floating-point
// errors. 100.10 + 100.30 is an IEEE 754 failure: float64 gives
// 200.39999999999998, not 200.40. See TestFloat64VsDecimalArithmetic for the proof.
func TestGetCartTotalMoneyMath(t *testing.T) {
	ts, cleanup := startServer(t)
	defer cleanup()

	sanitizerID := seedProduct(t, ts, "Hand Sanitizer", 10010, 10)
	maskID      := seedProduct(t, ts, "Face Mask",      10030, 10)

	token := signupAndLogin(t, ts, "buyer@example.com", "correct-password")
	post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: sanitizerID, Quantity: 1}, token, http.StatusOK, nil)
	post(t, ts, "/cart/items", cartapi.AddItemRequest{ProductID: maskID,      Quantity: 1}, token, http.StatusOK, nil)

	resp := &cartapi.CartTotalResponse{}
	get(t, ts, "/cart", token, http.StatusOK, resp)

	if resp.GetTotal() != "200.40" {
		t.Fatalf("expected total %q, got %q", "200.40", resp.GetTotal())
	}
}

// TestFloat64VsDecimalArithmetic isolates the same cart arithmetic used in
// TestGetCartTotalMoneyMath and shows what a float64-based total would return.
func TestFloat64VsDecimalArithmetic(t *testing.T) {
	// Simulate computeTotal with float64: convert paise to float, multiply by qty,
	// accumulate. 100.10 and 100.30 have no exact binary representation, so their
	// float64 sum is 200.39999999999998, not 200.40.
	sanitizer := float64(10010) / 100 // 100.10
	mask      := float64(10030) / 100 // 100.30
	float64Total := sanitizer*1 + mask*1
	if float64Total == 200.40 {
		t.Fatal("float64 unexpectedly returned the correct value")
	}
	if got := strconv.FormatFloat(float64Total, 'f', -1, 64); got != "200.39999999999998" {
		t.Fatalf("float64 total = %s, want 200.39999999999998", got)
	}

	// money.Amount stores 100.10 as 10010×10⁻² and 100.30 as 10030×10⁻²; Add is
	// integer addition (10010+10030=20040), yielding exactly "200.40".
	a, _ := money.NewAmount("INR", 10010, 2)
	b, _ := money.NewAmount("INR", 10030, 2)
	total, _ := a.Add(b)
	if got := total.Decimal().String(); got != "200.40" {
		t.Fatalf("decimal total = %s, want 200.40", got)
	}
}
