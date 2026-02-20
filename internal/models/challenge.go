package models

import (
	"gorm.io/gorm"
)

type Challenge struct {
	gorm.Model
	Title       string
	Description string
	Points      int
	Flag        string
}
