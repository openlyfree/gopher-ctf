package main

import (
	"gopher-ctf/internal"
	"gopher-ctf/internal/models"
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
		return
	}

	router := gin.Default()
	router.Use(sessions.Sessions("gopher_session", cookie.NewStore([]byte("really funky secret key that probably shouldn't be hardcoded"))))
	internal.RegisterRoutes(db, router)
	log.Fatal(router.Run("127.0.0.1:8080"))
}
