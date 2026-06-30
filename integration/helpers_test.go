package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saiteja/ecommerce/auth"
	authapi "github.com/saiteja/ecommerce/auth/api_models"
	authdao "github.com/saiteja/ecommerce/auth/dao"
	"github.com/saiteja/ecommerce/cart"
	cartdao "github.com/saiteja/ecommerce/cart/dao"
	"github.com/saiteja/ecommerce/cmd/httpserver"
	"github.com/saiteja/ecommerce/product"
	productdao "github.com/saiteja/ecommerce/product/dao"
)

type testServer struct {
	*httpserver.Server
	HTTP *httptest.Server
}

func startServer(t *testing.T) (*testServer, func()) {
	t.Helper()
	productStore := productdao.NewInMemoryProductStore()
	inventoryStore := productdao.NewInMemoryInventoryStore()
	credentialsStore := authdao.NewInMemoryUserCredentialsStore()
	sessionStore := authdao.NewInMemorySessionStore()
	authService := auth.NewService(credentialsStore, sessionStore)
	productService := product.NewService(productStore, inventoryStore)
	cartService := cart.NewService(cartdao.NewInMemoryCartStore(), productService)
	appServer := httpserver.New(authService, productService, cartService, productStore, inventoryStore)
	ts := &testServer{
		Server: appServer,
		HTTP:   httptest.NewServer(appServer),
	}
	return ts, ts.HTTP.Close
}

func signupAndLogin(t *testing.T, ts *testServer, email, password string) string {
	t.Helper()
	post(t, ts, "/signup", authapi.SignupRequest{Email: email, Password: password}, "", http.StatusCreated, nil)
	var resp authapi.LoginResponse
	post(t, ts, "/login", authapi.LoginRequest{Email: email, Password: password}, "", http.StatusOK, &resp)
	return resp.Token
}

func get(t *testing.T, ts *testServer, path, token string, wantStatus int, dst any) {
	t.Helper()
	doRequest(t, http.MethodGet, ts, path, nil, token, wantStatus, dst)
}

func del(t *testing.T, ts *testServer, path, token string, wantStatus int, dst any) {
	t.Helper()
	doRequest(t, http.MethodDelete, ts, path, nil, token, wantStatus, dst)
}

func patch(t *testing.T, ts *testServer, path string, body any, token string, wantStatus int, dst any) {
	t.Helper()
	doRequest(t, http.MethodPatch, ts, path, body, token, wantStatus, dst)
}

func post(t *testing.T, ts *testServer, path string, body any, token string, wantStatus int, dst any) {
	t.Helper()
	doRequest(t, http.MethodPost, ts, path, body, token, wantStatus, dst)
}

func doRequest(t *testing.T, method string, ts *testServer, path string, body any, token string, wantStatus int, dst any) {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode request body: %v", err)
		}
	}

	var req *http.Request
	var err error
	if body != nil {
		req, err = http.NewRequest(method, ts.HTTP.URL+path, &buf)
	} else {
		req, err = http.NewRequest(method, ts.HTTP.URL+path, nil)
	}
	if err != nil {
		t.Fatalf("create %s request for %s: %v", method, path, err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != wantStatus {
		t.Fatalf("%s %s: expected status %d, got %d", method, path, wantStatus, resp.StatusCode)
	}
	if dst != nil {
		if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
			t.Fatalf("decode response from %s %s: %v", method, path, err)
		}
	}
}
