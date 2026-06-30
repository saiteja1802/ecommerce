package models

type Cart struct {
	UserID string
	Items  []*CartItem
}

func (c *Cart) GetUserID() string {
	if c != nil {
		return c.UserID
	}
	return ""
}

func (c *Cart) GetItems() []*CartItem {
	if c != nil {
		return c.Items
	}
	return nil
}
