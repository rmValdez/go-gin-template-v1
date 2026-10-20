package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rmValdez/go-gin-template-v1/internal/app/config"
	"github.com/rmValdez/go-gin-template-v1/internal/infrastructure/database"
	"github.com/rmValdez/go-gin-template-v1/internal/infrastructure/redis"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/middleware"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/models"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/routes"
)

func main() {
	cfg := config.Load()

	// Connect to Redis
	redis.Connect(cfg.RedisURL)

	// Connect to Database (exits if unavailable)
	database.Connect(cfg.DatabaseURL)

	// Run auto-migrations (Prisma migrate equivalent)
	if err := database.DB.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Session{},
	); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	gin.SetMode(cfg.GinMode)

	// Use gin.New() instead of Default() so we can control our own middleware stack
	router := gin.New()
	
	// Add Gin's default logger
	router.Use(gin.Logger())

	// Add our custom middlewares
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.RequestID())
	router.Use(middleware.BodySizeLimit(2 << 20)) // 2 MB max body size
	router.Use(middleware.Cors(cfg.AllowedOrigins))
	
	// Apply a global rate limit (e.g., 100 requests per minute per IP)
	router.Use(middleware.RateLimiter("global", 100, time.Minute))

	if len(cfg.TrustedProxies) > 0 {
		if err := router.SetTrustedProxies(cfg.TrustedProxies); err != nil {
			log.Printf("Warning: failed to set trusted proxies: %v", err)
		}
	} else {
		_ = router.SetTrustedProxies(nil)
	}

	routes.Register(router, cfg)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if redis.Client != nil {
		if err := redis.Client.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v\n", err)
		} else {
			log.Println("Redis connection closed")
		}
	}

	if err := database.Close(); err != nil {
		log.Printf("Error closing database connection: %v\n", err)
	} else {
		log.Println("Database connection closed")
	}

	log.Println("Server exiting")
}