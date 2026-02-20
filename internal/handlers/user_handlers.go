package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/openlyfree/gopher-ctf/internal/models"
	"github.com/openlyfree/gopher-ctf/internal/shared"
	"github.com/openlyfree/gopher-ctf/ui/pages"
	"golang.org/x/crypto/bcrypt"
)

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.String(http.StatusInternalServerError, "Session error")
		return
	}
	c.Header("HX-Redirect", "/")
	c.Status(http.StatusOK)
}

func LoginHandler(c *gin.Context) {
	var user models.User
	if err := shared.DB.Where("username = ?", c.PostForm("username")).First(&user).Error; err != nil {
		c.String(http.StatusOK, "Invalid credentials")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.PostForm("password"))); err != nil {
		c.String(http.StatusOK, "Invalid credentials")
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("username", user.Username)
	if err := session.Save(); err != nil {
		c.String(http.StatusInternalServerError, "Session error")
		return
	}
	c.Header("HX-Redirect", "/")
	c.Status(http.StatusOK)
}

func RegisterHandler(c *gin.Context) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), 14)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating account")
		return
	}

	if err := shared.DB.Create(&models.User{Username: c.PostForm("username"), Password: string(hashedPassword)}).Error; err != nil {
		c.String(200, "Username already taken")
		return
	}

	c.Header("HX-Redirect", "/login")
	c.Status(http.StatusOK)
}

func ScoreboardHandler(c *gin.Context) {
	var users []models.User
	if err := shared.DB.Order("score desc").Find(&users).Error; err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	if c.GetHeader("HX-Request") == "true" {
		Render(c, 200, pages.Scoreboard(users))
	} else {
		Render(c, 200, pages.ScoreboardPage(users))
	}
}

func RegisterPageHandler(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		Render(c, 200, pages.Signup())
	} else {
		Render(c, 200, pages.SignupPage())
	}
}

func LoginPageHandler(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		Render(c, 200, pages.Login())
	} else {
		Render(c, 200, pages.LoginPage())
	}
}
