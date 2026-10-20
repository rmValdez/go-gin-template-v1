package security

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrPasswordTooLong is returned when a password exceeds bcrypt's 72-byte limit
	ErrPasswordTooLong = errors.New("password cannot exceed 72 bytes")

	// DummyHash is used during login when a user is not found to prevent timing attacks
	DummyHash string
)

const DefaultBcryptCost = 12

func init() {
	hash, err := bcrypt.GenerateFromPassword([]byte("timing-mitigation-dummy-password"), DefaultBcryptCost)
	if err == nil {
		DummyHash = string(hash)
	}
}

// HashPassword creates a bcrypt hash with cost 12 and validates length
func HashPassword(password string) (string, error) {
	if len([]byte(password)) > 72 {
		return "", ErrPasswordTooLong
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultBcryptCost)
	return string(bytes), err
}

// CheckPassword compares a plaintext password with a hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
