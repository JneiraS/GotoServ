package api

import (
	"net/http"
	"os"
	"time"

	cts "github.com/JneiraS/GotoServ/internal/constants"
	"github.com/JneiraS/GotoServ/pkg/utils"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.Use(rateLimitMiddleware(30, time.Minute))
	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 64*1024)
		c.Next()
	})
	secret := os.Getenv("SECRET_KEY")

	// Endpoint pour récupérer les assignments au format JSON
	router.GET("/:totp/assignments", func(c *gin.Context) {
		code := c.Param("totp")
		r, err := utils.ValidateTOTP(secret, code)
		if err != nil || !r {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.File(cts.AssignmentsJSON)

	})
	// Endpoint de santé pour vérifier que le serveur fonctionne
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Endpoint pour ajouter ou mettre à jour un assignment
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

		if err := utils.UpdateOrAddCSVRecord(cts.AssignmentsCSV, req.Agent, req.Scope, req.Keywords); err != nil {
			c.JSON(500, gin.H{"error": "failed to update CSV"})
			return
		}
		c.JSON(200, gin.H{"status": "CSV updated"})
	})

	// Endpoint pour mettre à jour les keywords d'un agent existant
	router.PATCH("/:totp/keywords", func(c *gin.Context) {
		code := c.Param("totp")
		r, err := utils.ValidateTOTP(secret, code)
		if err != nil || !r {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		var req struct {
			Agent    string `form:"agent" json:"agent" binding:"required"`
			Keywords string `form:"keywords" json:"keywords" binding:"required"`
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(400, gin.H{"error": "missing fields"})
			return
		}

		found, err := utils.UpdateKeywordsForAgent(cts.AssignmentsCSV, req.Agent, req.Keywords)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to update keywords"})
			return
		}
		if !found {
			c.JSON(404, gin.H{"error": "agent not found"})
			return
		}

		c.JSON(200, gin.H{"status": "keywords updated"})
	})

	return router
}
