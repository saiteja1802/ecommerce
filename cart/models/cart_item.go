package models

type CartItem struct {
	ProductID string
	Quantity  int
}

func (c *CartItem) GetProductID() string {
	if c != nil {
		return c.ProductID
	}
	return ""
}

func (c *CartItem) GetQuantity() int {
	if c != nil {
		return c.Quantity
	}
	return 0
}
