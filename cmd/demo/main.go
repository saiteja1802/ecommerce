package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	authapi "github.com/saiteja/ecommerce/auth/api_models"
	cartapi "github.com/saiteja/ecommerce/cart/api_models"
	productapi "github.com/saiteja/ecommerce/product/api_models"
)

var base string

func main() {
	base = os.Getenv("BASE_URL")
	if base == "" {
		base = "http://localhost:8080"
	}

	fmt.Println("=== ecommerce demo ===")

	fmt.Println("\n1. Signing up...")
	var signup authapi.SignupResponse
	post("/signup", authapi.SignupRequest{Email: "demo@example.com", Password: "demo1234"}, "", &signup)
	fmt.Printf("   user_id: %s\n", signup.UserID)

	fmt.Println("\n2. Logging in...")
	var login authapi.LoginResponse
	post("/login", authapi.LoginRequest{Email: "demo@example.com", Password: "demo1234"}, "", &login)
	token := login.Token
	fmt.Printf("   token: %s\n", token)

	fmt.Println("\n3. Product catalogue...")
	var catalog productapi.GetProductsCatalogResponse
	get("/products?page=1&page_size=10", token, &catalog)
	for _, p := range catalog.Products {
		fmt.Printf("   %-15s INR %s  (id: %s)\n", p.Name, p.Price, p.ID)
	}
	if len(catalog.Products) < 2 {
		log.Fatal("expected at least 2 products")
	}
	p1, p2 := catalog.Products[0].ID, catalog.Products[1].ID

	fmt.Println("\n4. Adding first product (x1)...")
	var r1 cartapi.AddItemResponse
	post("/cart/items", cartapi.AddItemRequest{ProductID: p1, Quantity: 1}, token, &r1)
	fmt.Printf("   cart total: INR %s\n", r1.Total)

	fmt.Println("\n5. Adding second product (x2)...")
	var r2 cartapi.AddItemResponse
	post("/cart/items", cartapi.AddItemRequest{ProductID: p2, Quantity: 2}, token, &r2)
	fmt.Printf("   cart total: INR %s\n", r2.Total)

	fmt.Println("\n6. Cart total (no coupon)...")
	var bare cartapi.GetCartTotalResponse
	get("/cart", token, &bare)
	fmt.Printf("   total: INR %s\n", bare.Total)

	fmt.Println("\n7. SAVE10 (10%% off, max 500.00)...")
	var s10 cartapi.GetCartTotalResponse
	get("/cart?coupon=SAVE10", token, &s10)
	fmt.Printf("   total: INR %s  discount: INR %s\n", s10.Total, s10.Discount)

	fmt.Println("\n8. WELCOME20 (20%% off, max 200.00)...")
	var w20 cartapi.GetCartTotalResponse
	get("/cart?coupon=WELCOME20", token, &w20)
	fmt.Printf("   total: INR %s  discount: INR %s\n", w20.Total, w20.Discount)

	fmt.Println("\n=== Done ===")
}

func post(path string, body any, token string, dst any) {
	doRequest(http.MethodPost, path, body, token, dst)
}

func get(path string, token string, dst any) {
	doRequest(http.MethodGet, path, nil, token, dst)
}

func doRequest(method, path string, body any, token string, dst any) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			log.Fatalf("encode request for %s %s: %v", method, path, err)
		}
	}

	req, err := http.NewRequest(method, base+path, &buf)
	if err != nil {
		log.Fatalf("create request for %s %s: %v", method, path, err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("%s %s: %v", method, path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Fatalf("%s %s: unexpected status %d", method, path, resp.StatusCode)
	}
	if dst != nil {
		if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
			log.Fatalf("decode response from %s %s: %v", method, path, err)
		}
	}
}
