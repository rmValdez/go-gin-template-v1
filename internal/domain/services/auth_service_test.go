package services_test

import (
	"testing"

	"github.com/rmValdez/go-gin-template-v1/internal/domain/apperrors"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/models"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/services"
)

// --- Mock User Repository ---
type mockUserRepository struct {
	users map[string]*models.User
}

func (m *mockUserRepository) Create(user *models.User) error {
	if _, exists := m.users[user.Email]; exists {
		return apperrors.ErrUserAlreadyExists
	}
	user.ID = uint(len(m.users) + 1)
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) FindByEmail(email string) (*models.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, nil // Typically GORM returns an error here, but for our mock nil, nil is fine if the service checks it
}

// --- Mock Session Repository ---
type mockSessionRepository struct {
	sessions map[string]*models.Session
}

func (m *mockSessionRepository) Create(session *models.Session) error {
	m.sessions[session.SessionToken] = session
	return nil
}

func (m *mockSessionRepository) FindByToken(token string) (*models.Session, error) {
	if session, exists := m.sessions[token]; exists {
		return session, nil
	}
	return nil, nil
}

func (m *mockSessionRepository) DeleteByToken(token string) error {
	delete(m.sessions, token)
	return nil
}

// --- The Actual Tests ---

func TestRegister(t *testing.T) {
	mockUserRepo := &mockUserRepository{users: make(map[string]*models.User)}
	mockSessionRepo := &mockSessionRepository{sessions: make(map[string]*models.Session)}
	
	// Inject the mock repositories into the actual service
	authService := services.NewAuthService(mockUserRepo, mockSessionRepo)

	t.Run("Successfully Register User", func(t *testing.T) {
		user, err := authService.Register("test@example.com", "securepassword123")
		
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		if user.Email != "test@example.com" {
			t.Errorf("Expected email to be test@example.com, got %s", user.Email)
		}
		
		if user.ID != 1 {
			t.Errorf("Expected ID to be 1, got %d", user.ID)
		}
	})

	t.Run("Fail on Duplicate Email", func(t *testing.T) {
		_, err := authService.Register("test@example.com", "anotherpassword")
		
		if err != apperrors.ErrUserAlreadyExists {
			t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
		}
	})
}

func TestLogin(t *testing.T) {
	mockUserRepo := &mockUserRepository{users: make(map[string]*models.User)}
	mockSessionRepo := &mockSessionRepository{sessions: make(map[string]*models.Session)}
	authService := services.NewAuthService(mockUserRepo, mockSessionRepo)

	// Setup a user directly in the mock to test login
	_, _ = authService.Register("login@example.com", "mypassword123")

	t.Run("Successfully Login", func(t *testing.T) {
		token, user, err := authService.Login("login@example.com", "mypassword123")
		
		if err != nil {
			t.Fatalf("Expected no error on login, got %v", err)
		}
		
		if token == "" {
			t.Errorf("Expected a session token, got empty string")
		}
		
		if user.Email != "login@example.com" {
			t.Errorf("Expected email to be login@example.com, got %s", user.Email)
		}

		// Verify session was stored in DB
		session, _ := mockSessionRepo.FindByToken(token)
		if session == nil {
			t.Errorf("Expected session to be stored in mock repository")
		}
	})

	t.Run("Fail with Incorrect Password", func(t *testing.T) {
		_, _, err := authService.Login("login@example.com", "wrongpassword")
		
		if err != apperrors.ErrInvalidCredentials {
			t.Errorf("Expected ErrInvalidCredentials, got %v", err)
		}
	})
}

func TestLogout(t *testing.T) {
	mockUserRepo := &mockUserRepository{users: make(map[string]*models.User)}
	mockSessionRepo := &mockSessionRepository{sessions: make(map[string]*models.Session)}
	authService := services.NewAuthService(mockUserRepo, mockSessionRepo)

	// Setup user and login to get a token
	_, _ = authService.Register("logout@example.com", "mypassword123")
	token, _, _ := authService.Login("logout@example.com", "mypassword123")

	t.Run("Successfully Logout", func(t *testing.T) {
		err := authService.Logout(token)
		
		if err != nil {
			t.Fatalf("Expected no error on logout, got %v", err)
		}

		// Verify session was deleted
		session, _ := mockSessionRepo.FindByToken(token)
		if session != nil {
			t.Errorf("Expected session to be deleted, but it still exists")
		}
	})
}
