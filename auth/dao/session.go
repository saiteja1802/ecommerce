package dao

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/saiteja/ecommerce/auth/models"
)

const sessionTTL = 15 * time.Minute

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

type InMemorySessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*models.Session
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]*models.Session),
	}
}

func (s *InMemorySessionStore) CreateSession(session *models.Session) (*models.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session.Token = uuid.New().String()
	if session.GetCreatedAt().IsZero() {
		session.CreatedAt = time.Now().UTC()
	}
	s.sessions[session.GetToken()] = session
	return session, nil
}

func (s *InMemorySessionStore) GetSession(token string) (*models.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[token]
	if !exists {
		return nil, ErrSessionNotFound
	}
	if time.Since(session.GetCreatedAt()) > sessionTTL {
		return nil, ErrSessionExpired
	}

	return session, nil
}
