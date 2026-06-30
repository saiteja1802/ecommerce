package api_models

import "time"

type GetProductDetailsRequest struct {
	ProductID string
}

func (r *GetProductDetailsRequest) GetProductID() string {
	if r != nil {
		return r.ProductID
	}
	return ""
}

type GetProductDetailsResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        string    `json:"price"`
	CurrencyCode string    `json:"currency_code"`
	CreatedAt    time.Time `json:"created_at"`
}

func (r *GetProductDetailsResponse) GetID() string {
	if r != nil {
		return r.ID
	}
	return ""
}

func (r *GetProductDetailsResponse) GetName() string {
	if r != nil {
		return r.Name
	}
	return ""
}

func (r *GetProductDetailsResponse) GetDescription() string {
	if r != nil {
		return r.Description
	}
	return ""
}

func (r *GetProductDetailsResponse) GetPrice() string {
	if r != nil {
		return r.Price
	}
	return ""
}

func (r *GetProductDetailsResponse) GetCurrencyCode() string {
	if r != nil {
		return r.CurrencyCode
	}
	return ""
}

func (r *GetProductDetailsResponse) GetCreatedAt() time.Time {
	if r != nil {
		return r.CreatedAt
	}
	return time.Time{}
}

type ProductSummary struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

func (p *ProductSummary) GetID() string {
	if p != nil {
		return p.ID
	}
	return ""
}

func (p *ProductSummary) GetName() string {
	if p != nil {
		return p.Name
	}
	return ""
}

func (p *ProductSummary) GetPrice() string {
	if p != nil {
		return p.Price
	}
	return ""
}

type GetInventoryRequest struct {
	ProductID string
}

func (r *GetInventoryRequest) GetProductID() string {
	if r != nil {
		return r.ProductID
	}
	return ""
}

type GetInventoryResponse struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (r *GetInventoryResponse) GetProductID() string {
	if r != nil {
		return r.ProductID
	}
	return ""
}

func (r *GetInventoryResponse) GetQuantity() int {
	if r != nil {
		return r.Quantity
	}
	return 0
}

type GetProductsCatalogRequest struct {
	Page     int
	PageSize int
}

func (r *GetProductsCatalogRequest) GetPage() int {
	if r != nil {
		return r.Page
	}
	return 0
}

func (r *GetProductsCatalogRequest) GetPageSize() int {
	if r != nil {
		return r.PageSize
	}
	return 0
}

type GetProductsCatalogResponse struct {
	Products []*ProductSummary `json:"products"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	TotalProducts int `json:"total_products"`
}

func (r *GetProductsCatalogResponse) GetProducts() []*ProductSummary {
	if r != nil {
		return r.Products
	}
	return nil
}

func (r *GetProductsCatalogResponse) GetPage() int {
	if r != nil {
		return r.Page
	}
	return 0
}

func (r *GetProductsCatalogResponse) GetPageSize() int {
	if r != nil {
		return r.PageSize
	}
	return 0
}

func (r *GetProductsCatalogResponse) GetTotalProducts() int {
	if r != nil {
		return r.TotalProducts
	}
	return 0
}
