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
			c.File(assignmentsFilePath + ".json")
		}

	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/:totp/add", func(c *gin.Context) {
		code := c.Param("totp")
		r, err := utils.ValidateTOTP(secret, code)
		if err != nil || !r {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		var req struct {
			Agent    string `form:"agent" json:"agent" binding:"required"`
			Scope    string `form:"scope" json:"scope" binding:"required"`
			Keywords string `form:"keywords" json:"keywords" binding:"required"`
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(400, gin.H{"error": "missing fields"})
			return
		}

		if err := utils.UpdateOrAddCSVRecord(assignmentsFilePath+".csv", req.Agent, req.Scope, req.Keywords); err != nil {
			c.JSON(500, gin.H{"error": "failed to update CSV"})
			return
		}
		utils.CreatJsonFromCsv()
		c.JSON(200, gin.H{"status": "CSV updated"})
	})

	return router
}
