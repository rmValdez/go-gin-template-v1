package repositories

import (
	"github.com/rmValdez/go-gin-template-v1/internal/domain/models"
)

type SessionRepository interface {
	Create(session *models.Session) error
	FindByToken(token string) (*models.Session, error)
	DeleteByToken(token string) error
}
