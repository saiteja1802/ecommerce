package dao

import "github.com/saiteja/ecommerce/auth/models"

type UserCredentialsDAO interface {
	CreateUserCredentials(credentials *models.UserCredentials) (*models.UserCredentials, error)
	GetCredentialsByEmail(email string) (*models.UserCredentials, error)
	GetCredentialsByUserID(userID string) (*models.UserCredentials, error)
}

type SessionDAO interface {
	CreateSession(session *models.Session) (*models.Session, error)
	GetSession(token string) (*models.Session, error)
}
