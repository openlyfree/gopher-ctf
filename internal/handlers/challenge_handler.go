package handlers

import (
	"gopher-ctf/internal/models"
	"gopher-ctf/internal/shared"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
