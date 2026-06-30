package models

type ProductInventory struct {
	ProductID string
	Quantity  int
}

func (p *ProductInventory) GetProductID() string {
	if p != nil {
		return p.ProductID
	}
	return ""
}

func (p *ProductInventory) GetQuantity() int {
	if p != nil {
		return p.Quantity
	}
	return 0
}
