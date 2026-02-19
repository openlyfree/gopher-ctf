package models

import "gorm.io/gorm"

type ChallengeCompletion struct {
	gorm.Model
	UserID      uint `gorm:"index;not null"`
	ChallengeID uint `gorm:"index;not null"`
}
