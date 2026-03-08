package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/openlyfree/gopher-ctf/internal/models"
	"github.com/openlyfree/gopher-ctf/internal/shared"
	"github.com/openlyfree/gopher-ctf/ui/components"
	"github.com/openlyfree/gopher-ctf/ui/pages"
	"gorm.io/gorm"
)

func SubmitFlagHandler(c *gin.Context) {
	userIDVal := sessions.Default(c).Get("user_id")
	if userIDVal == nil {
		c.String(200, "Login required")
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.String(200, "Session error")
		return
	}

	challengeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		c.String(200, "Invalid challenge ID")
		return
	}

	var challenge models.Challenge
	if err := shared.DB.First(&challenge, challengeID).Error; err != nil {
		_ = c.Error(err)
		c.String(200, "Challenge not found")
		return
	}

	submittedFlag := c.PostForm("flag")
	if submittedFlag != challenge.Flag {
		c.String(200, "Incorrect flag")
		return
	}

	var count int64
	if err := shared.DB.Model(&models.ChallengeCompletion{}).Where("user_id = ? AND challenge_id = ?", userID, challengeID).Count(&count).Error; err != nil {
		_ = c.Error(err)
		c.String(200, "Database error")

		return
	}
	if count > 0 {
		c.String(200, "Already completed!")
		return
	}

	tx := shared.DB.Begin()

	if err := tx.Create(&models.ChallengeCompletion{
		UserID:      userID,
		ChallengeID: uint(challengeID),
	}).Error; err != nil {
		tx.Rollback()
		_ = c.Error(err)
		c.String(200, "Database error")
		return
	}

	if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("score", gorm.Expr("score + ?", challenge.Points)).Error; err != nil {
		tx.Rollback()
		_ = c.Error(err)
		c.String(200, "Database error")
		return
	}

	if err := tx.Commit().Error; err != nil {
		_ = c.Error(err)
		c.String(200, "Database error")
		return
	}

	c.String(200, "Correct! +%d points", challenge.Points)
}

func GetChallenge(id string) (models.Challenge, error) {
	tempID, err := strconv.Atoi(id)
	if err != nil {
		return models.Challenge{}, err
	}
	var challenge models.Challenge
	if err := shared.DB.First(&challenge, tempID).Error; err != nil {
		return models.Challenge{}, err
	}
	return challenge, nil
}

func ChallengeIndividualHandler(c *gin.Context) {
	challenge, err := GetChallenge(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		c.Status(200)
		return
	}
	if c.GetHeader("HX-Request") == "true" {
		Render(c, 200, components.ChallengeCard(challenge))
	} else {
		Render(c, 200, pages.ChallengeDetail(challenge))
	}
}

func Render(c *gin.Context, status int, cmp templ.Component) {
	c.Status(status)
	_ = cmp.Render(c.Request.Context(), c.Writer)
}

func ChallengeIndexHandler(c *gin.Context) {
	var challenges []models.Challenge
	if err := shared.DB.Find(&challenges).Error; err != nil {
		_ = c.Error(err)
		c.String(200, "Database error")
		return
	}
	if c.GetHeader("HX-Request") == "true" {
		Render(c, 200, pages.Challenges(challenges))
	} else {
		Render(c, 200, pages.ChallengesPage(challenges))
	}
}

func CreateChallengeHandler(c *gin.Context) {
	name := c.PostForm("challenge_name")
	description := c.PostForm("challenge_description")
	pointsStr := c.PostForm("challenge_points")
	flag := c.PostForm("challenge_flag")

	if name == "" || description == "" || pointsStr == "" {
		c.String(200, "All fields are required")
		return
	}

	points, err := strconv.Atoi(pointsStr)
	if err != nil {
		_ = c.Error(err)
		c.String(200, "Invalid points value")
		return
	}

	challenge := models.Challenge{
		Title:       name,
		Description: description,
		Points:      points,
		Flag:        flag,
	}

	if err := shared.DB.Create(&challenge).Error; err != nil {
		_ = c.Error(err)
		c.String(200, "Database error")
		return
	}

	file, err := c.FormFile("challenge_files")
	if err == nil {
		challengeDir := filepath.Join("ctf", "challenges", fmt.Sprintf("%d", challenge.ID))
		if err := os.MkdirAll(challengeDir, 0o755); err != nil {
			_ = c.Error(err)
			c.String(200, "Challenge created but failed to create directory")
			return
		}
		dst := filepath.Join(challengeDir, file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			_ = c.Error(err)
			c.String(200, "Challenge created but failed to save file")
			return
		}
	}

	c.String(200, "Challenge '%s' created! (%d pts)", name, points)
}
