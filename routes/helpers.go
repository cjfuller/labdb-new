package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"labdb.org/labdb/models"
)

// /:model/:id/next
func nextRoute(cls string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		redirectID := models.NextID(cls, id)
		c.Redirect(307, fmt.Sprintf("/%s/%s", cls, redirectID))
	}
}

// /:model/:id/previous
func previousRoute(cls string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		redirectID := models.PrevID(cls, id)
		c.Redirect(307, fmt.Sprintf("/%s/%s", cls, redirectID))
	}
}
