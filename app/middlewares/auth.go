package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mdozairq/e-voting-be/app/helpers"
)

// VoterTokenAuthMiddleware verifies the voter token
func VoterTokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		voterToken := c.Request.Header.Get("Authorization")
		if voterToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		_, err := helpers.ValidateToken(voterToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Next()

	}
}

// CandidateTokenAuthMiddleware verifies the candidate token
func CandidateTokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		candidateToken := c.Request.Header.Get("Authorization")
		if candidateToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		_, err := helpers.ValidateCandidateToken(candidateToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminTokenAuthMiddleware verifies the admin token
func AdminTokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminToken := c.Request.Header.Get("Authorization")
		if adminToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		_, err := helpers.AdminValidateToken(adminToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Next()
	}
}