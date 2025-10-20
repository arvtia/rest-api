package middleware

import (
	"net/http"
	"strings"

	"github.com/arvtia/rest-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		switch claims.Role {
		case "admin":
			c.Set("adminID", claims.AdminID)
			c.Set("adminEmail", claims.Email)
			c.Set("role", "admin")
		case "user":
			c.Set("userID", claims.UserID)
			c.Set("userEmail", claims.Email)
			c.Set("role", "user")
		default:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unknown role"})
			return
		}

		c.Next()
	}
}
