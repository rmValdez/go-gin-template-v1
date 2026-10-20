package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/apperrors"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/services"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/response"
)

type CookieConfig struct {
	Domain   string
	Secure   bool
	SameSite http.SameSite
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,max=72"`
}

type AuthHandler struct {
	authService services.AuthService
	cookieCfg   CookieConfig
}

func NewAuthHandler(authService services.AuthService, cookieCfg CookieConfig) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cookieCfg:   cookieCfg,
	}
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

	// Set HttpOnly, Secure, SameSite session cookie
	c.SetSameSite(h.cookieCfg.SameSite)
	c.SetCookie("session_token", token, 3600*24*7, "/", h.cookieCfg.Domain, h.cookieCfg.Secure, true)

	response.Success(c, http.StatusOK, gin.H{
		"message": "Successfully logged in",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("userID")

	response.Success(c, http.StatusOK, gin.H{
		"message": "You are authenticated!",
		"userID":  userID,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := c.Cookie("session_token")
	if err == nil && token != "" {
		_ = h.authService.Logout(token)
	}

	// Clear the cookie
	c.SetSameSite(h.cookieCfg.SameSite)
	c.SetCookie("session_token", "", -1, "/", h.cookieCfg.Domain, h.cookieCfg.Secure, true)

	response.Success(c, http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}
