package api

import "github.com/gin-gonic/gin"

func NewRouter(assignmentsFilePath string) *gin.Engine {
	router := gin.Default()

	router.GET("/assignments", func(c *gin.Context) {
		c.File(assignmentsFilePath)
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
