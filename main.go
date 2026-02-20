package main

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/openlyfree/gopher-ctf/internal"
	"github.com/openlyfree/gopher-ctf/internal/shared"
)

func main() {
	shared.DB = shared.InitDB()

	router := gin.Default()
	router.Use(sessions.Sessions("gopher_session", cookie.NewStore([]byte("really funky secret key that probably shouldn't be hardcoded"))))
	internal.RegisterRoutes(router)
	log.Fatal(router.Run(":8080"))
}
