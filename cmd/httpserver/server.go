package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/saiteja/ecommerce/auth"
	authapi "github.com/saiteja/ecommerce/auth/api_models"
	"github.com/saiteja/ecommerce/cart"
	cartapi "github.com/saiteja/ecommerce/cart/api_models"
	cartdao "github.com/saiteja/ecommerce/cart/dao"
	"github.com/saiteja/ecommerce/pkg/logger"
	"github.com/saiteja/ecommerce/product"
	productapi "github.com/saiteja/ecommerce/product/api_models"
	productdao "github.com/saiteja/ecommerce/product/dao"
)

type Server struct {
	auth           *auth.Service
	product        *product.Service
	cart           *cart.Service
	mux            *http.ServeMux
	ProductStore   *productdao.InMemoryProductStore
	InventoryStore *productdao.InMemoryInventoryStore
	CouponStore    *cartdao.InMemoryCouponStore
}

type errorResponse struct {
	Error string `json:"error"`
}

func New(authService *auth.Service, productService *product.Service, cartService *cart.Service, productStore *productdao.InMemoryProductStore, inventoryStore *productdao.InMemoryInventoryStore, couponStore *cartdao.InMemoryCouponStore) *Server {
	logger.Init()
	s := &Server{
		auth:           authService,
		product:        productService,
		cart:           cartService,
		mux:            http.NewServeMux(),
		ProductStore:   productStore,
		InventoryStore: inventoryStore,
		CouponStore:    couponStore,
	}

	s.mux.HandleFunc("POST /signup", s.handleSignup)
	s.mux.HandleFunc("POST /login", s.handleLogin)
	s.mux.HandleFunc("POST /authenticate", s.handleAuthenticate)
	s.mux.HandleFunc("GET /products", s.handleGetProducts)
	s.mux.HandleFunc("GET /products/{id}", s.handleGetProductDetails)
	s.mux.HandleFunc("GET /products/{id}/inventory", s.handleGetProductInventory)
	s.mux.HandleFunc("POST /cart/items", s.handleAddToCart)
	s.mux.HandleFunc("DELETE /cart/items/{productID}", s.handleRemoveFromCart)
	s.mux.HandleFunc("PATCH /cart/items/{productID}", s.handleUpdateCartItem)
	s.mux.HandleFunc("GET /cart", s.handleGetCartTotal)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// authenticate extracts and validates the Bearer token, returning the userID.
func (s *Server) authenticate(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", auth.ErrInvalidToken
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	resp, err := s.auth.Authenticate(r.Context(), &authapi.AuthenticateRequest{Token: token})
	if err != nil {
		return "", err
	}
	return resp.GetUserID(), nil
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

func (s *Server) handleGetProductDetails(w http.ResponseWriter, r *http.Request) {
	if _, err := s.authenticate(r); err != nil {
		writeServiceError(w, err)
		return
	}

	resp, err := s.product.GetProductDetails(r.Context(), &productapi.GetProductDetailsRequest{
		ProductID: r.PathValue("id"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleGetProductInventory(w http.ResponseWriter, r *http.Request) {
	if _, err := s.authenticate(r); err != nil {
		writeServiceError(w, err)
		return
	}

	resp, err := s.product.GetInventory(r.Context(), &productapi.GetInventoryRequest{
		ProductID: r.PathValue("id"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	if _, err := s.authenticate(r); err != nil {
		writeServiceError(w, err)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	resp, err := s.product.GetProductsCatalog(r.Context(), &productapi.GetProductsCatalogRequest{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleAddToCart(w http.ResponseWriter, r *http.Request) {
	userID, err := s.authenticate(r)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	var req cartapi.AddItemRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	req.UserID = userID

	resp, err := s.cart.AddItem(r.Context(), &req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleRemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID, err := s.authenticate(r)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	resp, err := s.cart.RemoveItem(r.Context(), &cartapi.RemoveItemRequest{
		UserID:    userID,
		ProductID: r.PathValue("productID"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleUpdateCartItem(w http.ResponseWriter, r *http.Request) {
	userID, err := s.authenticate(r)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	var req cartapi.UpdateQuantityRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	req.UserID = userID
	req.ProductID = r.PathValue("productID")

	resp, err := s.cart.UpdateQuantity(r.Context(), &req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleGetCartTotal(w http.ResponseWriter, r *http.Request) {
	userID, err := s.authenticate(r)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	resp, err := s.cart.GetCartTotal(r.Context(), &cartapi.GetCartTotalRequest{
		UserID:     userID,
		CouponName: r.URL.Query().Get("coupon"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
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
	case errors.Is(err, auth.ErrInvalidToken), errors.Is(err, auth.ErrSessionExpired):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, cart.ErrInvalidQuantity):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, cart.ErrInsufficientStock):
		return http.StatusUnprocessableEntity, err.Error()
	case errors.Is(err, cart.ErrItemNotFound), errors.Is(err, cart.ErrProductNotFound),
		errors.Is(err, cart.ErrCouponNotFound), errors.Is(err, product.ErrProductNotFound):
		return http.StatusNotFound, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
