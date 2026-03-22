package api

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func StartServer(router *gin.Engine) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := "0.0.0.0:" + port

	log.Printf("server started on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
