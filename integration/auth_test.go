package integration_test

import (
	"net/http"
	"testing"

	authapi "github.com/saiteja/ecommerce/auth/api_models"
)

func TestSignupLoginAndAuthenticate(t *testing.T) {
	ts, cleanup := startServer(t)
	defer cleanup()

	var signupResp authapi.SignupResponse
	post(t, ts, "/signup", authapi.SignupRequest{
		Email:    "Buyer@Example.com",
		Password: "correct-password",
	}, "", http.StatusCreated, &signupResp)

	var loginResp authapi.LoginResponse
	post(t, ts, "/login", authapi.LoginRequest{
		Email:    "buyer@example.com",
		Password: "correct-password",
	}, "", http.StatusOK, &loginResp)

	if loginResp.Token == "" {
		t.Fatal("expected token in login response")
	}
	if loginResp.UserID != signupResp.UserID {
		t.Fatalf("login user id %q does not match signup user id %q", loginResp.UserID, signupResp.UserID)
	}

	var authResp authapi.AuthenticateResponse
	post(t, ts, "/authenticate", authapi.AuthenticateRequest{
		Token: loginResp.Token,
	}, "", http.StatusOK, &authResp)

	if authResp.UserID != signupResp.UserID {
		t.Fatalf("authenticated user id %q does not match signup user id %q", authResp.UserID, signupResp.UserID)
	}
}

func TestLoginRejectsInvalidPassword(t *testing.T) {
	ts, cleanup := startServer(t)
	defer cleanup()

	post(t, ts, "/signup", authapi.SignupRequest{
		Email:    "buyer@example.com",
		Password: "correct-password",
	}, "", http.StatusCreated, nil)

	post(t, ts, "/login", authapi.LoginRequest{
		Email:    "buyer@example.com",
		Password: "wrong-password",
	}, "", http.StatusUnauthorized, nil)
}

func TestDuplicateSignupIsRejected(t *testing.T) {
	ts, cleanup := startServer(t)
	defer cleanup()

	post(t, ts, "/signup", authapi.SignupRequest{
		Email:    "buyer@example.com",
		Password: "correct-password",
	}, "", http.StatusCreated, nil)

	post(t, ts, "/signup", authapi.SignupRequest{
		Email:    " BUYER@example.com ",
		Password: "another-password",
	}, "", http.StatusConflict, nil)
}

func TestAuthenticateRejectsInvalidToken(t *testing.T) {
	ts, cleanup := startServer(t)
	defer cleanup()

	post(t, ts, "/authenticate", authapi.AuthenticateRequest{
		Token: "not-a-real-token",
	}, "", http.StatusUnauthorized, nil)
}
