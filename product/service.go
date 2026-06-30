package product

import (
	"context"
	"errors"

	"github.com/saiteja/ecommerce/pkg/logger"
	productapi "github.com/saiteja/ecommerce/product/api_models"
	"github.com/saiteja/ecommerce/product/dao"
)

var ErrProductNotFound = errors.New("product not found")

type Service struct {
	productDao   dao.ProductDAO
	inventoryDao dao.InventoryDAO
}

func NewService(productDao dao.ProductDAO, inventoryDao dao.InventoryDAO) *Service {
	return &Service{productDao: productDao, inventoryDao: inventoryDao}
}

func (s *Service) GetProductDetails(ctx context.Context, req *productapi.GetProductDetailsRequest) (*productapi.GetProductDetailsResponse, error) {
	p, err := s.productDao.GetProductByID(req.GetProductID())
	if errors.Is(err, dao.ErrProductNotFound) {
		return nil, ErrProductNotFound
	}
	if err != nil {
		logger.L.Error("failed to get product by id", "error", err)
		return nil, err
	}

	price := p.GetPrice()
	return &productapi.GetProductDetailsResponse{
		ID:           p.GetID(),
		Name:         p.GetName(),
		Description:  p.GetDescription(),
		Price:        price.Decimal().String(),
		CurrencyCode: price.Curr().Code(),
		CreatedAt:    p.GetCreatedAt(),
	}, nil
}


func (s *Service) GetInventory(ctx context.Context, req *productapi.GetInventoryRequest) (*productapi.GetInventoryResponse, error) {
	inv, err := s.inventoryDao.GetInventory(req.GetProductID())
	if err != nil {
		if errors.Is(err, dao.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		logger.L.Error("failed to get inventory", "error", err)
		return nil, err
	}
	return &productapi.GetInventoryResponse{
		ProductID: inv.GetProductID(),
		Quantity:  inv.GetQuantity(),
	}, nil
}

func (s *Service) GetProductsCatalog(ctx context.Context, req *productapi.GetProductsCatalogRequest) (*productapi.GetProductsCatalogResponse, error) {
	page := req.GetPage()
	if page < 1 {
		page = 1
	}
	pageSize := req.GetPageSize()
	if pageSize < 1 {
		pageSize = 10
	}

	products, total, err := s.productDao.GetProducts(page, pageSize)
	if err != nil {
		logger.L.Error("failed to get products", "error", err)
		return nil, err
	}

	summaries := make([]*productapi.ProductSummary, len(products))
	for i, p := range products {
		summaries[i] = &productapi.ProductSummary{ID: p.GetID(), Name: p.GetName(), Price: p.GetPrice().Decimal().String()}
	}

	return &productapi.GetProductsCatalogResponse{
		Products:      summaries,
		Page:          page,
		PageSize:      pageSize,
		TotalProducts: total,
	}, nil
}
