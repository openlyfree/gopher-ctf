package main

import (
	"gopher-ctf/internal"
	"gopher-ctf/internal/models"
	"gopher-ctf/internal/shared"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, _ := gorm.Open(sqlite.Open("gopher-ctf.db"), &gorm.Config{})
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Challenge{}); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&models.ChallengeCompletion{}); err != nil {
		log.Fatal(err)
	}
	shared.DB = db

	router := gin.Default()
	router.Use(sessions.Sessions("gopher_session", cookie.NewStore([]byte("really funky secret key that probably shouldn't be hardcoded"))))
	internal.RegisterRoutes(router)
	log.Fatal(router.Run("127.0.0.1:8080"))
}
