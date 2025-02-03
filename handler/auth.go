package handler

import (
	"log/slog"
	"net/http"
	"safe-ollama/config"
	"safe-ollama/middleware"
	"safe-ollama/model"
	"safe-ollama/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthHandler(router *gin.Engine, db *gorm.DB) {
	r := router.Group("/api/auth")
	r.POST("/login", loginHandler(db))
	r.POST("/logout", logoutHandler(db))
	r.GET("/validateToken", middleware.LoginAuth(), validateTokenHandler(db))
}

func loginHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user model.User
		if err := db.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}

		hashedPassword := utils.EncryptPassword(loginRequest.Password, user.Salt)
		if hashedPassword != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}

		// 生成JWT令牌
		expirationTime := time.Now().Add(60 * time.Minute)
		claims := &model.JwtPayload{
			UserId:  user.ID,
			Role:    user.Role,
			Expires: expirationTime.Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.ToMapClaims())
		tokenString, err := token.SignedString(config.JwtSecret)
		if err != nil {
			slog.Error("Generate token failed", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		// 返回令牌给客户端
		c.JSON(http.StatusOK, gin.H{"userId": user.ID, "username": user.Username, "role": user.Role, "token": tokenString, "expires": expirationTime})
	}
}

func logoutHandler(_ *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	}
}

func validateTokenHandler(_ *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "token is valid"})
	}
}
