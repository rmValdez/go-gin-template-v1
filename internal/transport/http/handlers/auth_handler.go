package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/apperrors"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/response"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/services"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUserAlreadyExists):
			response.Error(c, http.StatusConflict, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidCredentials):
			response.Error(c, http.StatusUnauthorized, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	// Set the session cookie
	// Secure should be true in production, false for local dev unless using HTTPS
	c.SetCookie("session_token", token, 3600*24*7, "/", "", false, true)

	response.Success(c, http.StatusOK, gin.H{
		"message": "Successfully logged in",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	// Retrieved from the AuthRequired middleware
	userID, _ := c.Get("userID")
	
	response.Success(c, http.StatusOK, gin.H{
		"message": "You are authenticated!",
		"userID":  userID,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := c.Cookie("session_token")
	if err == nil && token != "" {
		// Attempt to delete the session from the database
		// We ignore the error because if it fails or is already deleted, 
		// we still want to clear the client's cookie.
		_ = h.authService.Logout(token)
	}

	// Clear the cookie by setting MaxAge to -1
	c.SetCookie("session_token", "", -1, "/", "", false, true)

	response.Success(c, http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}
