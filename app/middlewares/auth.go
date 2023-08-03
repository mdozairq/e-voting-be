package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mdozairq/e-voting-be/app/helpers"
)

// AdminAuthMiddleware is a custom middleware to validate the admin token
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminAuthToken := c.GetHeader("Authorization")

		// Replace "YOUR_ADMIN_AUTH_TOKEN" with your actual admin authentication token
		if adminAuthToken != "YOUR_ADMIN_AUTH_TOKEN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Continue to the next middleware or API handler
		c.Next()
	}
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		_, err := helpers.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		// c.Set("email", claims.Email)
		// c.Set("first_name", claims.First_name)
		// c.Set("last_name", claims.Last_name)
		// c.Set("uid", claims.Uid)

		c.Next()
	}
}