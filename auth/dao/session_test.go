package dao

import (
	"errors"
	"testing"
	"time"

	"github.com/saiteja/ecommerce/auth/models"
)

func TestInMemorySessionStore_CreateSession(t *testing.T) {
	fixed := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	var store *InMemorySessionStore

	tests := []struct {
		name  string
		input *models.Session
		check func(t *testing.T, got *models.Session)
	}{
		{
			name:  "generates non-empty token",
			input: &models.Session{UserID: "user-1"},
			check: func(t *testing.T, got *models.Session) {
				if got.GetToken() == "" {
					t.Fatal("expected non-empty token")
				}
			},
		},
		{
			name:  "sets CreatedAt when zero",
			input: &models.Session{UserID: "user-1"},
			check: func(t *testing.T, got *models.Session) {
				if got.GetCreatedAt().IsZero() {
					t.Fatal("expected CreatedAt to be set")
				}
			},
		},
		{
			name:  "preserves non-zero CreatedAt",
			input: &models.Session{UserID: "user-1", CreatedAt: fixed},
			check: func(t *testing.T, got *models.Session) {
				if !got.GetCreatedAt().Equal(fixed) {
					t.Fatalf("expected CreatedAt %v, got %v", fixed, got.GetCreatedAt())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store = NewInMemorySessionStore()
			got, err := store.CreateSession(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			tt.check(t, got)
		})
	}
}

func TestInMemorySessionStore_GetSession(t *testing.T) {
	var store *InMemorySessionStore

	tests := []struct {
		name       string
		setup      func() string // returns the token to look up
		wantErr    error
		wantUserID string
	}{
		{
			name: "valid non-expired session",
			setup: func() string {
				s, _ := store.CreateSession(&models.Session{UserID: "user-1"})
				return s.GetToken()
			},
			wantUserID: "user-1",
		},
		{
			name:    "unknown token returns ErrSessionNotFound",
			setup:   func() string { return "not-a-real-token" },
			wantErr: ErrSessionNotFound,
		},
		{
			name:    "empty token returns ErrSessionNotFound",
			setup:   func() string { return "" },
			wantErr: ErrSessionNotFound,
		},
		{
			name: "expired session returns ErrSessionExpired",
			setup: func() string {
				// back-date CreatedAt past the TTL to simulate an expired session;
				// CreateSession preserves non-zero CreatedAt values
				expired := time.Now().Add(-(sessionTTL + time.Second))
				s, _ := store.CreateSession(&models.Session{
					UserID:    "user-1",
					CreatedAt: expired,
				})
				return s.GetToken()
			},
			wantErr: ErrSessionExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store = NewInMemorySessionStore()
			token := tt.setup()

			session, err := store.GetSession(token)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr == nil {
				if session == nil {
					t.Fatal("expected session, got nil")
				}
				if session.GetUserID() != tt.wantUserID {
					t.Fatalf("expected UserID %q, got %q", tt.wantUserID, session.GetUserID())
				}
			}
		})
	}
}
