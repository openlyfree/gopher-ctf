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

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("secret_candy_vault_access")
		if err != nil || cookie != "levraiglooby26" { //yes its hardcoded fuck you
			c.AbortWithStatus(404)
			return
		}
		c.Next()
	}
}
