package handlers

import (
	"encoding/json"
	"net/http"
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
		c.String(http.StatusUnauthorized, "Login required")
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.String(http.StatusInternalServerError, "Session error")
		return
	}

	challengeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid challenge ID")
		return
	}

	var challenge models.Challenge
	if err := shared.DB.First(&challenge, challengeID).Error; err != nil {
		c.String(http.StatusNotFound, "Challenge not found")
		return
	}

	submittedFlag := c.PostForm("flag")
	if submittedFlag != challenge.Flag {
		c.String(http.StatusOK, "Incorrect flag")
		return
	}

	var count int64
	if err := shared.DB.Model(&models.ChallengeCompletion{}).Where("user_id = ? AND challenge_id = ?", userID, challengeID).Count(&count).Error; err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	if count > 0 {
		c.String(http.StatusOK, "Already completed!")
		return
	}

	tx := shared.DB.Begin()

	if err := tx.Create(&models.ChallengeCompletion{
		UserID:      userID,
		ChallengeID: uint(challengeID),
	}).Error; err != nil {
		tx.Rollback()
		c.String(http.StatusInternalServerError, "Database error")
		return
	}

	if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("score", gorm.Expr("score + ?", challenge.Points)).Error; err != nil {
		tx.Rollback()
		c.String(http.StatusInternalServerError, "Database error")
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}

	c.String(http.StatusOK, "Correct! +%d points", challenge.Points)
}

func GetChallenge(id string) (models.Challenge, error) {
	tempId, err := strconv.Atoi(id)
	if err != nil {
		return models.Challenge{}, err
	}
	var challenge models.Challenge
	if err := shared.DB.First(&challenge, tempId).Error; err != nil {
		return models.Challenge{}, err
	}
	return challenge, nil
}

func ChallengeIndividualHandler(c *gin.Context) {
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
}

func Render(c *gin.Context, status int, cmp templ.Component) {
	c.Status(status)
	_ = cmp.Render(c.Request.Context(), c.Writer)
}

func ChallengeIndexHandler(c *gin.Context) {
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
