package apperrors

import "errors"

// Centralized application errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrTokenGeneration    = errors.New("failed to generate authentication token")
	ErrInternalServer     = errors.New("internal server error")
)
