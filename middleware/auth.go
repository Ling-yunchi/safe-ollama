package middleware

import (
	"log/slog"
	"net/http"
	"safe-ollama/config"
	"safe-ollama/model"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Token")
		if tokenHeader == "" {
			slog.Info("[LoginAuth] Token is empty")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login first"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
			return config.JwtSecret, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		} else {
			var jwtPayload model.JwtPayload
			jwtPayload.FromMapClaims(claims)
			c.Set("claims", jwtPayload)
			c.Next()
		}
	}
}

func RoleAuth(role []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login First"})
			c.Abort()
			return
		}
		userRole := claims.(model.JwtPayload).Role
		success := false
		for _, r := range role {
			if r == userRole {
				success = true
				break
			}
		}
		if !success {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func OllamaAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Extract token from "Bearer token"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Verify token in the database
		var ollamaToken model.OllamaToken
		if err := db.Where("token = ?", token).First(&ollamaToken).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("ollamaToken", ollamaToken)

		c.Next()
	}
}
