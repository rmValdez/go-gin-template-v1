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

	// Connect to Database
	database.Connect(cfg.DatabaseURL)
	if database.DB != nil {
		// Run auto-migrations (Prisma migrate equivalent)
		database.DB.AutoMigrate(
			&models.User{},
			&models.Account{},
			&models.Session{},
		)
	}

	gin.SetMode(cfg.GinMode)

	// Use gin.New() instead of Default() so we can control our own middleware stack
	router := gin.New()
	
	// Add Gin's default logger
	router.Use(gin.Logger())

	// Add our custom middlewares
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.RequestID())
	router.Use(middleware.Cors())
	
	// Apply a global rate limit (e.g., 100 requests per minute per IP)
	router.Use(middleware.RateLimiter(100, time.Minute))

	_ = router.SetTrustedProxies(nil)

	routes.Register(router)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
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

	log.Println("Server exiting")
}