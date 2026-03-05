package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openlyfree/gopher-ctf/internal/shared"
	"github.com/openlyfree/gopher-ctf/ui/pages"
)

func AdminLoginHandler(c *gin.Context) {
	if c.PostForm("password") == shared.Config.Password {
		token, err := shared.EncryptAdminToken(shared.Config.Password)
		if err != nil {
			_ = c.Error(err)
			c.String(http.StatusInternalServerError, "Auth setup error")
			return
		}
		c.Status(http.StatusOK)
		c.SetCookie("secret_candy_vault_access", token, 3153600000, "/", "", false, true)
		c.Header("HX-Redirect", "/secret-candy-vault")
	} else {
		_ = c.Error(errors.New("incorrect password"))
		c.String(http.StatusOK, "Incorrect password")
	}
}

func AdminIndexHandler(c *gin.Context) { Render(c, 200, pages.Admin()) }

func AdminLoginPageHandler(c *gin.Context) { Render(c, 200, pages.AdminLogin()) }

func ConfigReloadHandler(c *gin.Context) {
	shared.InitConfig()
	c.String(http.StatusOK, "Config reloaded")
}
