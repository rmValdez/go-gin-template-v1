package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/response"
)

func Health(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{
		"status": "ok",
	})
}