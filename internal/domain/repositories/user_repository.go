package repositories

import (
	"github.com/rmValdez/go-gin-template-v1/internal/domain/models"
)

// UserRepository is the interface that services will depend on.
// This decouples the service from the underlying database implementation.
type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
}
