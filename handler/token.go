package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"safe-ollama/middleware"
	"safe-ollama/model"
	"safe-ollama/utils"
	"time"
)

func OllamaTokenHandler(router *gin.Engine, db *gorm.DB) {
	r := router.Group("/api/token", middleware.LoginAuth(), middleware.RoleAuth([]string{model.USER_ROLE, model.ADMIN_ROLE}))
	r.GET("/", getOllamaToken(db))
	r.POST("/", createOllamaToken(db))
	r.DELETE("/:tokenId", deleteOllamaToken(db))
}

type TokenResult struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
}

func getOllamaToken(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login First"})
			return
		}
		userId := claims.(model.JwtPayload).UserId

		var tokens []TokenResult
		if err := db.Model(&model.OllamaToken{}).Where("user_id = ?", userId).Find(&tokens).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tokens"})
			return
		}

		c.JSON(http.StatusOK, tokens)
	}
}

type CreateBean struct {
	Name string `json:"name"`
}

func createOllamaToken(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login First"})
			return
		}
		userId := claims.(model.JwtPayload).UserId
		var createBean CreateBean
		err := c.BindJSON(&createBean)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		name := createBean.Name
		if len(name) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name cannot be empty"})
			return
		}

		// 生成随机字符串作为token
		token := utils.GenerateToken(32)

		// 创建新的OllamaToken实例
		newToken := model.OllamaToken{
			Token:  token,
			Name:   name,
			UserId: userId,
		}

		// 保存到数据库
		if err := db.Create(&newToken).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"msg": "success"})
	}
}

func deleteOllamaToken(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login First"})
			return
		}
		userId := claims.(model.JwtPayload).UserId
		tokenId := c.Param("tokenId")

		// 获取要删除的token
		var token model.OllamaToken
		if err := db.Where("user_id = ? AND id = ?", userId, tokenId).First(&token).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete token"})
			return
		}

		if err := db.Delete(&token).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Token deleted successfully"})
	}
}
