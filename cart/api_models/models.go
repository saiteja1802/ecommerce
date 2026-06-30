package api_models

type AddItemRequest struct {
	UserID    string `json:"-"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (r *AddItemRequest) GetUserID() string {
	if r != nil {
		return r.UserID
	}
	return ""
}

func (r *AddItemRequest) GetProductID() string {
	if r != nil {
		return r.ProductID
	}
	return ""
}

func (r *AddItemRequest) GetQuantity() int {
	if r != nil {
		return r.Quantity
	}
	return 0
}

type AddItemResponse struct {
	Items []*CartItemSummary `json:"items"`
	Total string             `json:"total"`
}

func (r *AddItemResponse) GetItems() []*CartItemSummary {
	if r != nil {
		return r.Items
	}
	return nil
}

func (r *AddItemResponse) GetTotal() string {
	if r != nil {
		return r.Total
	}
	return ""
}

type RemoveItemRequest struct {
	UserID    string
	ProductID string
}

func (r *RemoveItemRequest) GetUserID() string {
	if r != nil {
		return r.UserID
	}
	return ""
}

func (r *RemoveItemRequest) GetProductID() string {
	if r != nil {
		return r.ProductID
	}
	return ""
}

type RemoveItemResponse struct {
	Items []*CartItemSummary `json:"items"`
	Total string             `json:"total"`
}

func (r *RemoveItemResponse) GetItems() []*CartItemSummary {
	if r != nil {
		return r.Items
	}
	return nil
}

func (r *RemoveItemResponse) GetTotal() string {
	if r != nil {
		return r.Total
	}
	return ""
}

type UpdateQuantityRequest struct {
	UserID    string `json:"-"`
	ProductID string `json:"-"`
	Quantity  int    `json:"quantity"`
}

func (r *UpdateQuantityRequest) GetUserID() string {
	if r != nil {
		return r.UserID
	}
	return ""
}

func (r *UpdateQuantityRequest) GetProductID() string {
	if r != nil {
		return r.ProductID
	}
	return ""
}

func (r *UpdateQuantityRequest) GetQuantity() int {
	if r != nil {
		return r.Quantity
	}
	return 0
}

type UpdateQuantityResponse struct {
	Items []*CartItemSummary `json:"items"`
	Total string             `json:"total"`
}

func (r *UpdateQuantityResponse) GetItems() []*CartItemSummary {
	if r != nil {
		return r.Items
	}
	return nil
}

func (r *UpdateQuantityResponse) GetTotal() string {
	if r != nil {
		return r.Total
	}
	return ""
}

type GetCartTotalRequest struct {
	UserID     string
	CouponName string
}

func (r *GetCartTotalRequest) GetUserID() string {
	if r != nil {
		return r.UserID
	}
	return ""
}

func (r *GetCartTotalRequest) GetCouponName() string {
	if r != nil {
		return r.CouponName
	}
	return ""
}

type CartItemSummary struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	UnitPrice   string `json:"unit_price"`
	Subtotal    string `json:"subtotal"`
}

func (r *CartItemSummary) GetProductID() string {
	if r != nil {
		return r.ProductID
	}
	return ""
}

func (r *CartItemSummary) GetProductName() string {
	if r != nil {
		return r.ProductName
	}
	return ""
}

func (r *CartItemSummary) GetQuantity() int {
	if r != nil {
		return r.Quantity
	}
	return 0
}

func (r *CartItemSummary) GetUnitPrice() string {
	if r != nil {
		return r.UnitPrice
	}
	return ""
}

func (r *CartItemSummary) GetSubtotal() string {
	if r != nil {
		return r.Subtotal
	}
	return ""
}

type GetCartTotalResponse struct {
	Items    []*CartItemSummary `json:"items"`
	Total    string             `json:"total"`
	Discount string             `json:"discount,omitempty"`
}

func (r *GetCartTotalResponse) GetItems() []*CartItemSummary {
	if r != nil {
		return r.Items
	}
	return nil
}

func (r *GetCartTotalResponse) GetTotal() string {
	if r != nil {
		return r.Total
	}
	return ""
}

func (r *GetCartTotalResponse) GetDiscount() string {
	if r != nil {
		return r.Discount
	}
	return ""
}
