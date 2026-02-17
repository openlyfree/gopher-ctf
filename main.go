package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ping", ping)
	log.Fatal(router.Run("127.0.0.1:8080"))
}
