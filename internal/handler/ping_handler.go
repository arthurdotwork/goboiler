package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		_ = ctx

		c.JSON(http.StatusOK, gin.H{})
	}
}
