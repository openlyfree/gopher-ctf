package internal

import (
	"context"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/openlyfree/gopher-ctf/internal/shared"
)

func UserLoader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if val := sessions.Default(c).Get("username"); val != nil {
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "username", val.(string)))
		}
		c.Next()
	}
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("secret_candy_vault_access")
		if err != nil || cookie != shared.Config.Password {
			c.AbortWithStatus(404)
			return
		}
		c.Next()
	}
}

func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}
