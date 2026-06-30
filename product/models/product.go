package models

import (
	"time"

	"github.com/govalues/money"
)

type Product struct {
	ID          string
	Name        string
	Description string
	Price       money.Amount
	CreatedAt   time.Time
}

func (p *Product) GetID() string {
	if p != nil {
		return p.ID
	}
	return ""
}

func (p *Product) GetName() string {
	if p != nil {
		return p.Name
	}
	return ""
}

func (p *Product) GetDescription() string {
	if p != nil {
		return p.Description
	}
	return ""
}

func (p *Product) GetPrice() money.Amount {
	if p != nil {
		return p.Price
	}
	return money.Amount{}
}

func (p *Product) GetCreatedAt() time.Time {
	if p != nil {
		return p.CreatedAt
	}
	return time.Time{}
}
