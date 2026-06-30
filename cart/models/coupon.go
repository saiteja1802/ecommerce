package models

import (
	"github.com/govalues/decimal"
	"github.com/govalues/money"
)

const InrScale = 2 // decimal places for INR (paise)

type Coupon struct {
	Name               string
	DiscountPercentage decimal.Decimal
	MaxDiscount        money.Amount
}

func (c *Coupon) GetName() string {
	if c != nil {
		return c.Name
	}
	return ""
}

func (c *Coupon) GetDiscountPercentage() decimal.Decimal {
	if c != nil {
		return c.DiscountPercentage
	}
	return decimal.Zero
}

func (c *Coupon) GetMaxDiscount() money.Amount {
	if c != nil {
		return c.MaxDiscount
	}
	zero, _ := money.NewAmount("INR", 0, InrScale)
	return zero
}
