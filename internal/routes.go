package internal

import (
	"gopher-ctf/internal/models"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"gopher-ctf/ui/pages"
)

var db *gorm.DB

func RegisterRoutes(database *gorm.DB, router *gin.Engine) {
	db = database
	router.Use(UserLoader())
	router.Static("/static", "./ui/static")
	router.GET("/", HandleIndex)
	router.GET("/login", func(c *gin.Context) { Render(c, 200, pages.Login()) })
	router.POST("/login", LoginHandler)
	router.GET("/register", func(c *gin.Context) { Render(c, 200, pages.Signup()) })
	router.POST("/register", RegisterHandler)
	router.POST("/logout", LogoutHandler)
}
func Render(c *gin.Context, status int, cmp templ.Component) {
	c.Status(status)
	_ = cmp.Render(c.Request.Context(), c.Writer)
}

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.String(http.StatusInternalServerError, "Failed to clear session")
		return
	}
	c.Header("HX-Redirect", "/")
	c.Status(http.StatusOK)
}
func HandleIndex(c *gin.Context) {
	Render(c, 200, pages.Index())
}

func LoginHandler(c *gin.Context) {
	var user models.User
	if err := db.Where("username = ?", c.PostForm("username")).First(&user).Error; err != nil {
		c.String(http.StatusUnauthorized, "Invalid credentials")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.PostForm("password"))); err != nil {
		c.String(http.StatusUnauthorized, "Invalid credentials")
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("username", user.Username)
	_ = session.Save()
	c.Header("HX-Redirect", "/dashboard")
	c.Status(http.StatusOK)
}

func RegisterHandler(c *gin.Context) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), 14)

	if err := db.Create(&models.User{Username: c.PostForm("username"), Password: string(hashedPassword)}).Error; err != nil {
		c.String(200, "Username already taken")
		return
	}

	c.Header("HX-Redirect", "/login")
	c.Status(http.StatusOK)
}
