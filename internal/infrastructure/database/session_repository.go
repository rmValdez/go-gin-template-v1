package database

import (
	"github.com/rmValdez/go-gin-template-v1/internal/domain/models"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/repositories"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) repositories.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *sessionRepository) FindByToken(token string) (*models.Session, error) {
	var session models.Session
	err := r.db.Where("session_token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) DeleteByToken(token string) error {
	return r.db.Where("session_token = ?", token).Delete(&models.Session{}).Error
}
