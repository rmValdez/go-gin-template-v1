package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rmValdez/go-gin-template-v1/internal/infrastructure/database"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/handlers"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/middleware"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/services"
)

func registerAuthRoutes(router *gin.RouterGroup) {
	// Dependency Injection: Wire up the layers!
	// We inject the infrastructure implementation into the domain interface
	userRepo := database.NewUserRepository(database.DB)
	sessionRepo := database.NewSessionRepository(database.DB)
	
	authService := services.NewAuthService(userRepo, sessionRepo)
	authHandler := handlers.NewAuthHandler(authService)

	auth := router.Group("/auth")
	
	// Public routes
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Protected routes using Auth middleware
	protected := auth.Group("/")
	protected.Use(middleware.AuthRequired(sessionRepo))
	{
		protected.GET("/me", authHandler.Me)
		protected.POST("/logout", authHandler.Logout)
	}
}
