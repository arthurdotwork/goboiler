package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/arthureichelberger/goboiler/internal/metrics"
	"github.com/gin-gonic/gin"
)

func InstrumentedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		elapsed := time.Since(start).Milliseconds()

		if c.Request.Method == http.MethodOptions {
			return
		}

		metrics.CountHTTPRequest(c.Request.Method, c.FullPath(), strconv.Itoa(c.Writer.Status()))
		metrics.ObserveHTTPRequestDuration(c.Request.Method, c.FullPath(), strconv.Itoa(c.Writer.Status()), float64(elapsed))
	}
}
