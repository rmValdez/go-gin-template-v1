package services

import (
	"strings"
	"time"

	"github.com/rmValdez/go-gin-template-v1/internal/domain/apperrors"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/models"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/repositories"
	"github.com/rmValdez/go-gin-template-v1/internal/infrastructure/security"
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
	// Normalize email to prevent whitespace and case-sensitive duplicate accounts
	email = strings.ToLower(strings.TrimSpace(email))

	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, apperrors.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		if err == security.ErrPasswordTooLong {
			return nil, apperrors.ErrInvalidCredentials
		}
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
	// Normalize email
	email = strings.ToLower(strings.TrimSpace(email))

	// 1. Fetch user through repository
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		// Run dummy password check to defend against timing-based user enumeration attacks
		_ = security.CheckPassword(password, security.DummyHash)
		return "", nil, apperrors.ErrInvalidCredentials
	}

	// 2. Verify password logic (business rule)
	if !security.CheckPassword(password, user.Password) {
		return "", nil, apperrors.ErrInvalidCredentials
	}

	// 3. Generate Secure Session Token
	rawToken, err := security.GenerateRandomToken(32)
	if err != nil {
		return "", nil, apperrors.ErrInternalServer
	}

	// 4. Save SHA-256 Hashed Session Token to Database
	hashedToken := security.HashToken(rawToken)
	session := &models.Session{
		SessionToken: hashedToken,
		UserID:       user.ID,
		ExpiresAt:    time.Now().Add(24 * time.Hour * 7), // 7 days
	}
	if err := s.sessionRepo.Create(session); err != nil {
		return "", nil, apperrors.ErrInternalServer
	}

	return rawToken, user, nil
}

func (s *authService) Logout(sessionToken string) error {
	hashedToken := security.HashToken(sessionToken)
	if err := s.sessionRepo.DeleteByToken(hashedToken); err != nil {
		return apperrors.ErrInternalServer
	}
	return nil
}
