package handler

import (
	"net/http"
	"safe-ollama/middleware"
	"safe-ollama/model"
	"safe-ollama/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserHandler(router *gin.Engine, db *gorm.DB) {
	r := router.Group("/api/user", middleware.LoginAuth(), middleware.RoleAuth([]string{model.ADMIN_ROLE}))
	r.GET("/:id", getUserInfo(db))
	r.POST("/", createUser(db))
	r.PUT("/:id", updateUser(db))
	r.DELETE("/:id", deleteUser(db))
	r.GET("/", getAllUser(db))
}

type UserBean struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResult struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func getUserInfo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user UserResult
		if err := db.Model(&model.User{}).First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func createUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userBean UserBean
		if err := c.ShouldBindJSON(&userBean); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		salt := utils.GenerateSalt()
		hashedPassword := utils.EncryptPassword(userBean.Password, salt)
		user := model.User{
			Username: userBean.Username,
			Password: hashedPassword,
			Salt:     salt,
			Role:     model.USER_ROLE,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

func updateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var userBean UserBean
		if err := c.ShouldBindJSON(&userBean); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user model.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if userBean.Password != "" {
			salt := utils.GenerateSalt()
			hashedPassword := utils.EncryptPassword(userBean.Password, salt)
			user.Password = hashedPassword
			user.Salt = salt
		}
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func deleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user model.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if id == "1" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can not delete admin user"})
			return
		}
		if err := db.Delete(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

func getAllUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []UserResult
		if err := db.Model(&model.User{}).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}
