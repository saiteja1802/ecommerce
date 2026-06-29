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
	"github.com/saiteja/ecommerce/cmd/httpserver"
	"github.com/saiteja/ecommerce/pkg/logger"
)

func TestSignupLoginAndAuthenticate(t *testing.T) {
	srv := startServer(t)
	defer srv.Close()

	var signupResp authapi.SignupResponse
	post(t, srv, "/signup", authapi.SignupRequest{
		Email:    "Buyer@Example.com",
		Password: "correct-password",
	}, http.StatusCreated, &signupResp)

	var loginResp authapi.LoginResponse
	post(t, srv, "/login", authapi.LoginRequest{
		Email:    "buyer@example.com",
		Password: "correct-password",
	}, http.StatusOK, &loginResp)

	if loginResp.Token == "" {
		t.Fatal("expected token in login response")
	}
	if loginResp.UserID != signupResp.UserID {
		t.Fatalf("login user id %q does not match signup user id %q", loginResp.UserID, signupResp.UserID)
	}

	var authResp authapi.AuthenticateResponse
	post(t, srv, "/authenticate", authapi.AuthenticateRequest{
		Token: loginResp.Token,
	}, http.StatusOK, &authResp)

	if authResp.UserID != signupResp.UserID {
		t.Fatalf("authenticated user id %q does not match signup user id %q", authResp.UserID, signupResp.UserID)
	}
}

func TestLoginRejectsInvalidPassword(t *testing.T) {
	srv := startServer(t)
	defer srv.Close()

	post(t, srv, "/signup", authapi.SignupRequest{
		Email:    "buyer@example.com",
		Password: "correct-password",
	}, http.StatusCreated, nil)

	post(t, srv, "/login", authapi.LoginRequest{
		Email:    "buyer@example.com",
		Password: "wrong-password",
	}, http.StatusUnauthorized, nil)
}

func TestDuplicateSignupIsRejected(t *testing.T) {
	srv := startServer(t)
	defer srv.Close()

	post(t, srv, "/signup", authapi.SignupRequest{
		Email:    "buyer@example.com",
		Password: "correct-password",
	}, http.StatusCreated, nil)

	post(t, srv, "/signup", authapi.SignupRequest{
		Email:    " BUYER@example.com ",
		Password: "another-password",
	}, http.StatusConflict, nil)
}

func TestAuthenticateRejectsInvalidToken(t *testing.T) {
	srv := startServer(t)
	defer srv.Close()

	post(t, srv, "/authenticate", authapi.AuthenticateRequest{
		Token: "not-a-real-token",
	}, http.StatusUnauthorized, nil)
}

func startServer(t *testing.T) *httptest.Server {
	t.Helper()
	logger.Init()
	credentialsStore := authdao.NewInMemoryUserCredentialsStore()
	sessionStore := authdao.NewInMemorySessionStore()
	authService := auth.NewService(credentialsStore, sessionStore)
	return httptest.NewServer(httpserver.New(authService))
}

func post(t *testing.T, srv *httptest.Server, path string, body any, wantStatus int, dst any) {
	t.Helper()

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		t.Fatalf("encode request body: %v", err)
	}

	resp, err := http.Post(srv.URL+path, "application/json", &buf)
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != wantStatus {
		t.Fatalf("POST %s: expected status %d, got %d", path, wantStatus, resp.StatusCode)
	}
	if dst != nil {
		if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
			t.Fatalf("decode response from POST %s: %v", path, err)
		}
	}
}
