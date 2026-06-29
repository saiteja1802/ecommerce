package main

import (
	"net/http"

	"github.com/saiteja/ecommerce/auth"
	authdao "github.com/saiteja/ecommerce/auth/dao"
	"github.com/saiteja/ecommerce/cmd/httpserver"
	"github.com/saiteja/ecommerce/pkg/logger"
)

func main() {
	credentialsStore := authdao.NewInMemoryUserCredentialsStore()
	sessionStore := authdao.NewInMemorySessionStore()
	authService := auth.NewService(credentialsStore, sessionStore)
	s := httpserver.New(authService)

	logger.L.Info("server started", "addr", ":8080")
	if err := http.ListenAndServe(":8080", s); err != nil {
		logger.L.Error("server stopped", "error", err)
	}
}
