package cart

import (
	"context"
	"errors"

	"github.com/govalues/decimal"
	"github.com/govalues/money"
	cartapi "github.com/saiteja/ecommerce/cart/api_models"
	cartdao "github.com/saiteja/ecommerce/cart/dao"
	cartmodels "github.com/saiteja/ecommerce/cart/models"
	"github.com/saiteja/ecommerce/pkg/logger"
	"github.com/saiteja/ecommerce/product"
	productapi "github.com/saiteja/ecommerce/product/api_models"
)

var (
	ErrInvalidQuantity   = errors.New("quantity must be at least 1")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrItemNotFound      = errors.New("item not found")
	ErrProductNotFound   = errors.New("product not found")
	ErrCouponNotFound    = errors.New("coupon not found")
)

type Service struct {
	cartDAO        cartdao.CartDAO
	couponDAO      cartdao.CouponDAO
	productService *product.Service
}

func NewService(cartDAO cartdao.CartDAO, couponDAO cartdao.CouponDAO, productService *product.Service) *Service {
	return &Service{
		cartDAO:        cartDAO,
		couponDAO:      couponDAO,
		productService: productService,
	}
}

// AddItem merges quantity into any existing cart entry for the same product.
func (s *Service) AddItem(ctx context.Context, req *cartapi.AddItemRequest) (*cartapi.AddItemResponse, error) {
	if req.GetQuantity() < 1 {
		return nil, ErrInvalidQuantity
	}

	inv, err := s.productService.GetInventory(ctx, &productapi.GetInventoryRequest{ProductID: req.GetProductID()})
	if errors.Is(err, product.ErrProductNotFound) {
		return nil, ErrProductNotFound
	}
	if err != nil {
		logger.L.Error("failed to get inventory", "productID", req.GetProductID(), "error", err)
		return nil, err
	}

	cart, err := s.cartDAO.GetCart(req.GetUserID())
	if err != nil && !errors.Is(err, cartdao.ErrCartNotFound) {
		logger.L.Error("failed to get cart", "error", err)
		return nil, err
	}
	newQty := req.GetQuantity()
	if cart != nil {
		for _, item := range cart.GetItems() {
			if item.GetProductID() == req.GetProductID() {
				newQty = item.GetQuantity() + req.GetQuantity()
				break
			}
		}
	}

	if inv.GetQuantity() < newQty {
		return nil, ErrInsufficientStock
	}

	if err := s.cartDAO.SetItem(req.GetUserID(), req.GetProductID(), newQty); err != nil {
		logger.L.Error("failed to set cart item", "error", err)
		return nil, err
	}

	items, total, err := s.computeTotal(ctx, req.GetUserID())
	if err != nil {
		return nil, err
	}
	return &cartapi.AddItemResponse{Items: items, Total: total}, nil
}

func (s *Service) RemoveItem(ctx context.Context, req *cartapi.RemoveItemRequest) (*cartapi.RemoveItemResponse, error) {
	err := s.cartDAO.RemoveItem(req.GetUserID(), req.GetProductID())
	if errors.Is(err, cartdao.ErrItemNotFound) {
		return nil, ErrItemNotFound
	}
	if err != nil {
		logger.L.Error("failed to remove cart item", "error", err)
		return nil, err
	}

	items, total, err := s.computeTotal(ctx, req.GetUserID())
	if err != nil {
		return nil, err
	}
	return &cartapi.RemoveItemResponse{Items: items, Total: total}, nil
}

// UpdateQuantity sets an absolute quantity for an item already in the cart.
// A quantity of 0 removes the item.
func (s *Service) UpdateQuantity(ctx context.Context, req *cartapi.UpdateQuantityRequest) (*cartapi.UpdateQuantityResponse, error) {
	if req.GetQuantity() < 0 {
		return nil, ErrInvalidQuantity
	}

	// Treat quantity 0 as an explicit remove.
	if req.GetQuantity() == 0 {
		if err := s.cartDAO.RemoveItem(req.GetUserID(), req.GetProductID()); err != nil {
			logger.L.Error("failed to remove cart item", "error", err)
			return nil, err
		}
		items, total, err := s.computeTotal(ctx, req.GetUserID())
		if err != nil {
			return nil, err
		}
		return &cartapi.UpdateQuantityResponse{Items: items, Total: total}, nil
	}

	cart, err := s.cartDAO.GetCart(req.GetUserID())
	if errors.Is(err, cartdao.ErrCartNotFound) {
		return nil, ErrItemNotFound
	}
	if err != nil {
		logger.L.Error("failed to get cart", "error", err)
		return nil, err
	}
	found := false
	for _, item := range cart.GetItems() {
		if item.GetProductID() == req.GetProductID() {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrItemNotFound
	}

	inv, err := s.productService.GetInventory(ctx, &productapi.GetInventoryRequest{ProductID: req.GetProductID()})
	if err != nil {
		logger.L.Error("failed to get inventory", "productID", req.GetProductID(), "error", err)
		return nil, err
	}
	if inv.GetQuantity() < req.GetQuantity() {
		return nil, ErrInsufficientStock
	}

	if err := s.cartDAO.SetItem(req.GetUserID(), req.GetProductID(), req.GetQuantity()); err != nil {
		logger.L.Error("failed to set cart item", "error", err)
		return nil, err
	}

	items, total, err := s.computeTotal(ctx, req.GetUserID())
	if err != nil {
		return nil, err
	}
	return &cartapi.UpdateQuantityResponse{Items: items, Total: total}, nil
}

// GetCartTotal computes the cart total. If a coupon name is provided, the
// discount is computed and subtracted; both the final total and the discount
// are returned separately in the response.
func (s *Service) GetCartTotal(ctx context.Context, req *cartapi.GetCartTotalRequest) (*cartapi.GetCartTotalResponse, error) {
	items, rawTotal, err := s.computeTotal(ctx, req.GetUserID())
	if err != nil {
		return nil, err
	}

	if req.GetCouponName() == "" {
		return &cartapi.GetCartTotalResponse{Items: items, Total: rawTotal}, nil
	}

	coupon, err := s.couponDAO.GetCoupon(req.GetCouponName())
	if errors.Is(err, cartdao.ErrCouponNotFound) {
		return nil, ErrCouponNotFound
	}
	if err != nil {
		logger.L.Error("failed to get coupon", "coupon", req.GetCouponName(), "error", err)
		return nil, err
	}

	totalAmt, err := money.ParseAmount("INR", rawTotal)
	if err != nil {
		logger.L.Error("failed to parse cart total", "error", err)
		return nil, err
	}

	discount, err := totalAmt.Mul(coupon.GetDiscountPercentage())
	if err != nil {
		logger.L.Error("failed to compute discount", "error", err)
		return nil, err
	}
	discount = discount.Round(cartmodels.InrScale)

	cmp, err := discount.Cmp(coupon.GetMaxDiscount())
	if err != nil {
		logger.L.Error("failed to compare discount to max discount", "error", err)
		return nil, err
	}
	if cmp > 0 {
		discount = coupon.GetMaxDiscount()
	}

	finalTotal, err := totalAmt.Sub(discount)
	if err != nil {
		logger.L.Error("failed to subtract discount", "error", err)
		return nil, err
	}

	return &cartapi.GetCartTotalResponse{
		Items:    items,
		Total:    finalTotal.Decimal().String(),
		Discount: discount.Decimal().String(),
	}, nil
}

func (s *Service) computeTotal(ctx context.Context, userID string) ([]*cartapi.CartItemSummary, string, error) {
	cart, err := s.cartDAO.GetCart(userID)
	if errors.Is(err, cartdao.ErrCartNotFound) {
		zero, _ := money.NewAmount("INR", 0, cartmodels.InrScale)
		return []*cartapi.CartItemSummary{}, zero.Decimal().String(), nil
	}
	if err != nil {
		logger.L.Error("failed to get cart", "error", err)
		return nil, "", err
	}

	total, _ := money.NewAmount("INR", 0, cartmodels.InrScale)
	summaries := make([]*cartapi.CartItemSummary, 0, len(cart.GetItems()))

	for _, item := range cart.GetItems() {
		p, err := s.productService.GetProductDetails(ctx, &productapi.GetProductDetailsRequest{ProductID: item.GetProductID()})
		if errors.Is(err, product.ErrProductNotFound) {
			return nil, "", ErrProductNotFound
		}
		if err != nil {
			logger.L.Error("failed to get product details", "productID", item.GetProductID(), "error", err)
			return nil, "", err
		}

		price, err := money.ParseAmount("INR", p.GetPrice())
		if err != nil {
			logger.L.Error("failed to parse product price", "price", p.GetPrice(), "error", err)
			return nil, "", err
		}

		qty, _ := decimal.New(int64(item.GetQuantity()), 0)
		subtotal, err := price.Mul(qty)
		if err != nil {
			logger.L.Error("failed to compute subtotal", "error", err)
			return nil, "", err
		}

		total, err = total.Add(subtotal)
		if err != nil {
			logger.L.Error("failed to accumulate total", "error", err)
			return nil, "", err
		}

		summaries = append(summaries, &cartapi.CartItemSummary{
			ProductID:   p.GetID(),
			ProductName: p.GetName(),
			Quantity:    item.GetQuantity(),
			UnitPrice:   price.Decimal().String(),
			Subtotal:    subtotal.Decimal().String(),
		})
	}

	return summaries, total.Decimal().String(), nil
}
