package handlers

import (
	"encoding/json"
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
	file, err := c.FormFile("challenge_file")
	if err != nil {
		_ = c.Error(err)
		c.String(200, "File upload error")
		return
	}
	openedFile, err := file.Open()
	if err != nil {
		_ = c.Error(err)
		c.String(200, "File open error")
		return
	}

	var newChallenges []models.Challenge
	if err := json.NewDecoder(openedFile).Decode(&newChallenges); err != nil {
		_ = c.Error(err)
		c.String(200, "Invalid JSON format")
		return
	}
	_ = openedFile.Close()
	for _, challenge := range newChallenges {
		if err := shared.DB.Create(&challenge).Error; err != nil {
			_ = c.Error(err)
			c.String(200, "Database error")
			return
		}
	}

	c.String(200, "Successfully uploaded %d challenges!", len(newChallenges))
}
