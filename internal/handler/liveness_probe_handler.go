package handler

import (
	"net/http"

	"github.com/arthureichelberger/goboiler/pkg/psql"
	"github.com/gin-gonic/gin"
)

func LivenessProbeHandler(db psql.Queryable) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := db.ExecContext(c.Request.Context(), "SELECT 1"); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
