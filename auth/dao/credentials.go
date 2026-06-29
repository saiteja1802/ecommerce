package dao

import (
	"errors"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/saiteja/ecommerce/auth/models"
)

var (
	ErrCredentialsNotFound = errors.New("credentials not found")
	ErrEmailExists         = errors.New("email already exists")
)

type InMemoryUserCredentialsStore struct {
	mu            sync.RWMutex
	usersByID     map[string]*models.UserCredentials
	userIDByEmail map[string]string
}

func NewInMemoryUserCredentialsStore() *InMemoryUserCredentialsStore {
	return &InMemoryUserCredentialsStore{
		usersByID:     make(map[string]*models.UserCredentials),
		userIDByEmail: make(map[string]string),
	}
}

func (s *InMemoryUserCredentialsStore) CreateUserCredentials(credentials *models.UserCredentials) (*models.UserCredentials, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.userIDByEmail[credentials.GetEmail()]; exists {
		return nil, ErrEmailExists
	}

	credentials.UserID = s.nextID()
	if credentials.GetCreatedAt().IsZero() {
		credentials.CreatedAt = time.Now().UTC()
	}

	s.usersByID[credentials.GetUserID()] = credentials
	s.userIDByEmail[credentials.GetEmail()] = credentials.GetUserID()

	return credentials, nil
}

func (s *InMemoryUserCredentialsStore) GetCredentialsByEmail(email string) (*models.UserCredentials, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, exists := s.userIDByEmail[email]
	if !exists {
		return nil, ErrCredentialsNotFound
	}

	return s.usersByID[userID], nil
}

func (s *InMemoryUserCredentialsStore) GetCredentialsByUserID(userID string) (*models.UserCredentials, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	credentials, exists := s.usersByID[userID]
	if !exists {
		return nil, ErrCredentialsNotFound
	}

	return credentials, nil
}

func (s *InMemoryUserCredentialsStore) nextID() string {
	return "USR" + ulid.Make().String()
}
