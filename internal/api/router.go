package api

import (
	"net/http"
	"os"
	"strings"
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
	router.GET("/assignments", func(c *gin.Context) {
		code := totpCodeFromHeader(c)
		r, err := utils.ValidateTOTP(secret, code)
		if err != nil || !r {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.File(cts.AssignmentsJSON)

	})
	// Endpoint de santé pour vérifier que le serveur fonctionne
	router.GET("/health", func(c *gin.Context) {
		code, err := utils.GenerateCurrentTOTP(secret)
		if err == nil {
			// Affiche le TOTP courant dans les logs à chaque requête /health
			println("TOTP actuel:", code)
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Endpoint pour ajouter ou mettre à jour un assignment
	router.POST("/add", func(c *gin.Context) {
		var req struct {
			TOTP     string `form:"totp" json:"totp"`
			Agent    string `form:"agent" json:"agent" binding:"required"`
			Scope    string `form:"scope" json:"scope" binding:"required"`
			Keywords string `form:"keywords" json:"keywords" binding:"required"`
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(400, gin.H{"error": "missing fields"})
			return
		}

		code := totpCodeFromRequest(c, req.TOTP)
		r, err := utils.ValidateTOTP(secret, code)
		if err != nil || !r {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		if err := utils.UpdateOrAddCSVRecord(cts.AssignmentsCSV, req.Agent, req.Scope, req.Keywords); err != nil {
			c.JSON(500, gin.H{"error": "failed to update CSV"})
			return
		}
		c.JSON(200, gin.H{"status": "CSV updated"})
	})

	// Endpoint pour mettre à jour les keywords d'un agent existant
	router.PATCH("/keywords", func(c *gin.Context) {
		var req struct {
			TOTP     string `form:"totp" json:"totp"`
			Agent    string `form:"agent" json:"agent" binding:"required"`
			Keywords string `form:"keywords" json:"keywords" binding:"required"`
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(400, gin.H{"error": "missing fields"})
			return
		}

		code := totpCodeFromRequest(c, req.TOTP)
		r, err := utils.ValidateTOTP(secret, code)
		if err != nil || !r {
			c.JSON(401, gin.H{"error": "unauthorized"})
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

func totpCodeFromHeader(c *gin.Context) string {
	return strings.TrimSpace(c.GetHeader("X-TOTP-Code"))
}

func totpCodeFromRequest(c *gin.Context, bodyCode string) string {
	if headerCode := totpCodeFromHeader(c); headerCode != "" {
		return headerCode
	}

	return strings.TrimSpace(bodyCode)
}
