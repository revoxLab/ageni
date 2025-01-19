package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PrePack() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Keys = make(map[string]interface{})
		c.Next()

		func() {
			if c.Writer.Written() {
				return
			}

			params := c.Keys
			if len(params) == 0 {
				return
			}
			c.JSON(http.StatusOK, params)
		}()
	}
}
