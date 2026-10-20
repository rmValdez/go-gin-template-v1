package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rmValdez/go-gin-template-v1/internal/app/config"
)

func Register(router *gin.Engine, cfg config.Config) {
	api := router.Group("/api/v1")

	registerHealthRoutes(api)
	registerAuthRoutes(api, cfg)
}
