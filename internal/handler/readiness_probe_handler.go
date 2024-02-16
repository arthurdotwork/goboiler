package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReadinessProbeHandler(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		select {
		case <-ctx.Done():
			c.JSON(http.StatusServiceUnavailable, gin.H{})
			return
		default:
			c.JSON(http.StatusOK, gin.H{})
		}
	}
}
