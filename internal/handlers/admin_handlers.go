package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openlyfree/gopher-ctf/ui/pages"
)

func AdminLoginHandler(c *gin.Context) {
	if c.PostForm("password") == "levraiglooby26" {
		c.Status(http.StatusOK)
		c.SetCookie("secret_candy_vault_access", "levraiglooby26", 3153600000, "/", "", false, true)
		c.Header("HX-Redirect", "/secret-candy-vault")
	} else {
		c.Error(errors.New("incorrect password"))
		c.String(http.StatusOK, "Incorrect password")
	}
}

func AdminIndexHandler(c *gin.Context) { Render(c, 200, pages.Admin()) }

func AdminLoginPageHandler(c *gin.Context) { Render(c, 200, pages.AdminLogin()) }
