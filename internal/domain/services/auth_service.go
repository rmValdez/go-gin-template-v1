package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/rmValdez/go-gin-template-v1/internal/domain/apperrors"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/models"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/repositories"
	"github.com/rmValdez/go-gin-template-v1/internal/infrastructure/security"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(email, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
	Logout(sessionToken string) error
}

type authService struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
}

func NewAuthService(userRepo repositories.UserRepository, sessionRepo repositories.SessionRepository) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *authService) Register(email, password string) (*models.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, apperrors.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	// Create new user
	user := &models.User{
		Email:    email,
		Password: hashedPassword,
	}

	// Save via repository
	if err := s.userRepo.Create(user); err != nil {
		return nil, apperrors.ErrInternalServer
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	// 1. Fetch user through the repository interface
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, apperrors.ErrInvalidCredentials
		}
		return "", nil, err
	}

	// 2. Verify password logic (business rule)
	if !security.CheckPassword(password, user.Password) {
		return "", nil, apperrors.ErrInvalidCredentials
	}

	// 3. Generate Secure Session Token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", nil, apperrors.ErrInternalServer
	}
	sessionToken := hex.EncodeToString(bytes)

	// 4. Save Session to Database
	session := &models.Session{
		SessionToken: sessionToken,
		UserID:       user.ID,
		ExpiresAt:    time.Now().Add(24 * time.Hour * 7), // 7 days
	}
	if err := s.sessionRepo.Create(session); err != nil {
		return "", nil, apperrors.ErrInternalServer
	}

	return sessionToken, user, nil
}

func (s *authService) Logout(sessionToken string) error {
	// Delete the session from the database
	if err := s.sessionRepo.DeleteByToken(sessionToken); err != nil {
		return apperrors.ErrInternalServer
	}
	return nil
}
