package models

import "time"

type UserCredentials struct {
	UserID       string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

func (c *UserCredentials) GetUserID() string {
	return c.UserID
}

func (c *UserCredentials) GetEmail() string {
	return c.Email
}

func (c *UserCredentials) GetPasswordHash() string {
	return c.PasswordHash
}

func (c *UserCredentials) GetCreatedAt() time.Time {
	return c.CreatedAt
}
