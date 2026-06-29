package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/mail"
	"strings"
	"time"

	authapi "github.com/saiteja/ecommerce/auth/api_models"
	"github.com/saiteja/ecommerce/auth/dao"
	"github.com/saiteja/ecommerce/auth/models"
	"github.com/saiteja/ecommerce/pkg/logger"
)

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrSessionExpired     = errors.New("session expired")
)

type Service struct {
	userCredentialsDao dao.UserCredentialsDAO
	sessionsDao        dao.SessionDAO
}

func NewService(userCredentialsDao dao.UserCredentialsDAO, sessionsDao dao.SessionDAO) *Service {
	return &Service{userCredentialsDao: userCredentialsDao, sessionsDao: sessionsDao}
}

func (s *Service) Signup(ctx context.Context, request *authapi.SignupRequest) (*authapi.SignupResponse, error) {
	email := normalizeEmail(request.GetEmail())
	if !validEmail(email) {
		return nil, ErrInvalidEmail
	}
	if len(request.GetPassword()) < 8 {
		return nil, ErrInvalidPassword
	}

	passwordHash, err := hashPassword(request.GetPassword())
	if err != nil {
		logger.L.Error("failed to hash password", "error", err)
		return nil, err
	}

	credentials, err := s.userCredentialsDao.CreateUserCredentials(&models.UserCredentials{
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	})
	if errors.Is(err, dao.ErrEmailExists) {
		return nil, ErrEmailAlreadyExists
	}
	if err != nil {
		logger.L.Error("failed to create user credentials", "error", err)
		return nil, err
	}

	return &authapi.SignupResponse{
		UserID:    credentials.GetUserID(),
		Email:     credentials.GetEmail(),
		CreatedAt: credentials.GetCreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *Service) Login(ctx context.Context, request *authapi.LoginRequest) (*authapi.LoginResponse, error) {
	credentials, err := s.userCredentialsDao.GetCredentialsByEmail(normalizeEmail(request.GetEmail()))
	if errors.Is(err, dao.ErrCredentialsNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		logger.L.Error("failed to get credentials by email", "error", err)
		return nil, err
	}

	if !checkPassword(credentials.GetPasswordHash(), request.GetPassword()) {
		return nil, ErrInvalidCredentials
	}

	session, err := s.sessionsDao.CreateSession(&models.Session{
		UserID: credentials.GetUserID(),
	})
	if err != nil {
		logger.L.Error("failed to create session", "error", err)
		return nil, err
	}

	return &authapi.LoginResponse{
		UserID: credentials.GetUserID(),
		Token:  session.GetToken(),
	}, nil
}

func (s *Service) Authenticate(ctx context.Context, request *authapi.AuthenticateRequest) (*authapi.AuthenticateResponse, error) {
	token := strings.TrimSpace(request.GetToken())
	if token == "" {
		return nil, ErrInvalidToken
	}

	session, err := s.sessionsDao.GetSession(token)
	if errors.Is(err, dao.ErrSessionNotFound) {
		return nil, ErrInvalidToken
	}
	if errors.Is(err, dao.ErrSessionExpired) {
		return nil, ErrSessionExpired
	}
	if err != nil {
		logger.L.Error("failed to get session", "error", err)
		return nil, err
	}

	return &authapi.AuthenticateResponse{UserID: session.GetUserID()}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	sum := sha256.Sum256(append(salt, []byte(password)...))
	return base64.RawURLEncoding.EncodeToString(salt) + "." + base64.RawURLEncoding.EncodeToString(sum[:]), nil
}

func checkPassword(encoded, password string) bool {
	parts := strings.Split(encoded, ".")
	if len(parts) != 2 {
		return false
	}

	salt, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	expected, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	actual := sha256.Sum256(append(salt, []byte(password)...))
	return subtle.ConstantTimeCompare(actual[:], expected) == 1
}

func validEmail(email string) bool {
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email
}
