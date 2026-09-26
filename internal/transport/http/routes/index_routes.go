package routes

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	api := router.Group("/api/v1")

	registerHealthRoutes(api)
	registerAuthRoutes(api)
}
