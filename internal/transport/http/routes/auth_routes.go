package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rmValdez/go-gin-template-v1/internal/app/config"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/services"
	"github.com/rmValdez/go-gin-template-v1/internal/infrastructure/database"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/handlers"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/middleware"
)

func registerAuthRoutes(router *gin.RouterGroup, cfg config.Config) {
	userRepo := database.NewUserRepository(database.DB)
	sessionRepo := database.NewSessionRepository(database.DB)

	authService := services.NewAuthService(userRepo, sessionRepo)

	sameSite := http.SameSiteLaxMode
	if strings.EqualFold(cfg.CookieSameSite, "strict") {
		sameSite = http.SameSiteStrictMode
	} else if strings.EqualFold(cfg.CookieSameSite, "none") {
		sameSite = http.SameSiteNoneMode
	}

	cookieCfg := handlers.CookieConfig{
		Domain:   cfg.CookieDomain,
		Secure:   cfg.CookieSecure,
		SameSite: sameSite,
	}

	authHandler := handlers.NewAuthHandler(authService, cookieCfg)

	auth := router.Group("/auth")

	// S4: Strict rate limit for authentication endpoints (10 attempts per minute per IP)
	authRateLimiter := middleware.RateLimiter("auth", 10, time.Minute)

	// Public routes
	auth.POST("/register", authRateLimiter, authHandler.Register)
	auth.POST("/login", authRateLimiter, authHandler.Login)

	// Protected routes using Auth middleware
	protected := auth.Group("/")
	protected.Use(middleware.AuthRequired(sessionRepo))
	{
		protected.GET("/me", authHandler.Me)
		protected.POST("/logout", authHandler.Logout)
	}
}
