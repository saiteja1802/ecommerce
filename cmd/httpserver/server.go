package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/saiteja/ecommerce/auth"
	authapi "github.com/saiteja/ecommerce/auth/api_models"
	"github.com/saiteja/ecommerce/pkg/logger"
)

type Server struct {
	auth *auth.Service
	mux  *http.ServeMux
}

type errorResponse struct {
	Error string `json:"error"`
}

func New(authService *auth.Service) *Server {
	logger.Init()
	s := &Server{
		auth: authService,
		mux:  http.NewServeMux(),
	}

	s.mux.HandleFunc("POST /signup", s.handleSignup)
	s.mux.HandleFunc("POST /login", s.handleLogin)
	s.mux.HandleFunc("POST /authenticate", s.handleAuthenticate)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	var request authapi.SignupRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	response, err := s.auth.Signup(r.Context(), &request)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, response)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var request authapi.LoginRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	response, err := s.auth.Login(r.Context(), &request)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleAuthenticate(w http.ResponseWriter, r *http.Request) {
	var request authapi.AuthenticateRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	response, err := s.auth.Authenticate(r.Context(), &request)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeServiceError(w http.ResponseWriter, err error) {
	status, message := statusCodeForError(err)
	writeError(w, status, message)
}

func statusCodeForError(err error) (int, string) {
	switch {
	case errors.Is(err, auth.ErrInvalidEmail), errors.Is(err, auth.ErrInvalidPassword):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, auth.ErrEmailAlreadyExists):
		return http.StatusConflict, err.Error()
	case errors.Is(err, auth.ErrInvalidCredentials):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, auth.ErrInvalidToken):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, auth.ErrSessionExpired):
		return http.StatusUnauthorized, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
