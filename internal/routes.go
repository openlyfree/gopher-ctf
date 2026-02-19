package internal

import (
	"encoding/json"
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
	router.Static("/static", "./ui/static")
	router.Use(UserLoader())

	router.GET("/", func(c *gin.Context) { Render(c, 200, pages.Index()) })
	router.GET("/login", func(c *gin.Context) { Render(c, 200, pages.Login()) })
	router.GET("/register", func(c *gin.Context) { Render(c, 200, pages.Signup()) })
	router.GET("/gimme-secret-candy-vault", func(c *gin.Context) { Render(c, 200, pages.AdminLogin()) })
	router.GET("/challenges", func(c *gin.Context) { Render(c, 200, pages.Challenges(db)) })

	router.POST("/login", LoginHandler)
	router.POST("/register", RegisterHandler)
	router.POST("/logout", LogoutHandler)
	router.POST("/gimme-secret-candy-vault", AdminLoginHandler)
	admin := router.Group("/secret-candy-vault")
	admin.Use(AdminAuth())

	admin.GET("/", func(c *gin.Context) { Render(c, 200, pages.Admin()) })

	admin.POST("/challenge", CreateChallengeHandler)

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
	c.Header("HX-Redirect", "/")
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

func CreateChallengeHandler(c *gin.Context) {
	file, _ := c.FormFile("challenge_file")
	openedFile, _ := file.Open()
	defer openedFile.Close()

	var newChallenges []models.Challenge
	if err := json.NewDecoder(openedFile).Decode(&newChallenges); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON format")
		return
	}

	for _, challenge := range newChallenges {
		db.Create(&challenge)
	}

	c.String(http.StatusOK, "Successfully uploaded %d challenges!", len(newChallenges))
}
func AdminLoginHandler(c *gin.Context) {
	if c.PostForm("password") == "levraiglooby26" {
		c.Status(http.StatusOK)
		c.SetCookie("secret_candy_vault_access", "levraiglooby26", 3153600000, "/", "", false, true)
		c.Header("HX-Redirect", "/secret-candy-vault")
	} else {
		c.String(http.StatusOK, "Incorrect password")
	}
}
