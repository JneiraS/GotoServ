package api

import (
	"os"

	"github.com/JneiraS/GotoServ/pkg/utils"
	"github.com/gin-gonic/gin"
)

func NewRouter(assignmentsFilePath string) *gin.Engine {
	router := gin.Default()
	secret := os.Getenv("SECRET_KEY")

	router.GET("/:totp/assignments", func(c *gin.Context) {
		code := c.Param("totp")
		r, err := utils.ValidateTOTP(secret, code)
		if err == nil && r == true {
			c.File(assignmentsFilePath)
		}

	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
