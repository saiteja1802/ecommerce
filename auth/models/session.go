package models

import "time"

type Session struct {
	// Token is set by the CreateSession method, callers need not populate this field
	Token     string
	UserID    string
	CreatedAt time.Time
}

func (s *Session) GetToken() string {
	return s.Token
}

func (s *Session) GetUserID() string {
	return s.UserID
}

func (s *Session) GetCreatedAt() time.Time {
	return s.CreatedAt
}
