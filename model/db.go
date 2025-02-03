package model

import (
	"errors"
	"safe-ollama/config"
	"safe-ollama/utils"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.DatabasePath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = InitModels(db)
	if err != nil {
		panic(err)
	}

	seedData(db)

	return db
}

func seedData(db *gorm.DB) {
	var user User
	adminUsername := config.GetStringWithDefault("admin.username", "admin")
	adminPassword := config.GetStringWithDefault("admin.password", "admin")

	if err := db.Where("username = ?", adminUsername).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			salt := utils.GenerateSalt()
			adminPassword = utils.EncryptPassword(adminPassword, salt)
			db.Create(&User{
				Username: adminUsername,
				Password: adminPassword,
				Salt:     salt,
				Role:     ADMIN_ROLE,
			})
		} else {
			panic(err)
		}
	}
}
