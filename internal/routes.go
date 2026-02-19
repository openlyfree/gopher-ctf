package internal

import (
	"encoding/json"
	"gopher-ctf/internal/handlers"
	"gopher-ctf/internal/models"
	"gopher-ctf/internal/shared"
	"gopher-ctf/ui/components"
	"net/http"
	"strconv"

	"gopher-ctf/ui/pages"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.Static("/static", "./ui/static")
	router.Use(UserLoader())

	router.GET("/", func(c *gin.Context) { Render(c, 200, pages.Index()) })
	router.GET("/login", func(c *gin.Context) {
		if c.GetHeader("HX-Request") == "true" {
			Render(c, 200, pages.Login())
		} else {
			Render(c, 200, pages.LoginPage())
		}
	})
	router.GET("/register", func(c *gin.Context) {
		if c.GetHeader("HX-Request") == "true" {
			Render(c, 200, pages.Signup())
		} else {
			Render(c, 200, pages.SignupPage())
		}
	})
	router.GET("/gimme-secret-candy-vault", func(c *gin.Context) { Render(c, 200, pages.AdminLogin()) })
	router.GET("/challenges", func(c *gin.Context) {
		var challenges []models.Challenge
		if err := shared.DB.Find(&challenges).Error; err != nil {
			c.String(http.StatusInternalServerError, "Database error")
			return
		}
		if c.GetHeader("HX-Request") == "true" {
			Render(c, 200, pages.Challenges(challenges))
		} else {
			Render(c, 200, pages.ChallengesPage(challenges))
		}
	})

	router.POST("/login", handlers.LoginHandler)
	router.POST("/register", handlers.RegisterHandler)
	router.POST("/logout", handlers.LogoutHandler)
	router.POST("/gimme-secret-candy-vault", handlers.AdminLoginHandler)

	challengeGroup := router.Group("/challenges")
	challengeGroup.GET("/:id", func(c *gin.Context) {
		challenge, err := GetChallenge(c.Param("id"))
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		if c.GetHeader("HX-Request") == "true" {
			Render(c, 200, components.ChallengeCard(challenge))
		} else {
			Render(c, 200, pages.ChallengeDetail(challenge))
		}
	})
	challengeGroup.POST("/:id/submit", handlers.SubmitFlagHandler)

	admin := router.Group("/secret-candy-vault")
	admin.Use(AdminAuth())
	admin.GET("/", func(c *gin.Context) { Render(c, 200, pages.Admin()) })
	admin.POST("/challenge", CreateChallengeHandler)

	router.GET("/scoreboard", func(c *gin.Context) {
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
	})

}
func Render(c *gin.Context, status int, cmp templ.Component) {
	c.Status(status)
	_ = cmp.Render(c.Request.Context(), c.Writer)
}

func CreateChallengeHandler(c *gin.Context) {
	file, err := c.FormFile("challenge_file")
	if err != nil {
		c.String(http.StatusBadRequest, "File upload error")
		return
	}
	openedFile, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "File open error")
		return
	}
	defer openedFile.Close()

	var newChallenges []models.Challenge
	if err := json.NewDecoder(openedFile).Decode(&newChallenges); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON format")
		return
	}

	for _, challenge := range newChallenges {
		if err := shared.DB.Create(&challenge).Error; err != nil {
			c.String(http.StatusInternalServerError, "Database error")
			return
		}
	}

	c.String(http.StatusOK, "Successfully uploaded %d challenges!", len(newChallenges))
}

func GetChallenge(id string) (models.Challenge, error) {
	tempId, err := strconv.Atoi(id)
	if err != nil {
		tempId = 1
	}
	var challenge models.Challenge
	if err := shared.DB.First(&challenge, tempId).Error; err != nil {
		return models.Challenge{}, err
	}
	return challenge, nil
}
