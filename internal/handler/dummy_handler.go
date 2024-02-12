package handler

import (
	"net/http"

	"github.com/arthureichelberger/goboiler/internal/service"
	"github.com/gin-gonic/gin"
)

func DummyHandler(dummyService service.DummyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		count, err := dummyService.Dummy(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"count": count})
	}
}
