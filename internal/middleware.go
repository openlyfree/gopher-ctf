package internal

import (
	"context"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func UserLoader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if val := sessions.Default(c).Get("username"); val != nil {
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "username", val.(string)))
		}
		c.Next()
	}
}
