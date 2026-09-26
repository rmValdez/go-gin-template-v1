package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/handlers"
)

func registerHealthRoutes(router *gin.RouterGroup) {
	router.GET("/health", handlers.Health)
}
