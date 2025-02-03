package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	ADMIN_ROLE = "admin"
	USER_ROLE  = "user"
)

type User struct {
	ID       uint   `gorm:"primaryKey; autoIncrement"`
	Username string `gorm:"not null; index:user_username_index"`
	Password string `gorm:"not null"`
	Salt     string `gorm:"not null"`
	Role     string `gorm:"not null"`
}

type OllamaToken struct {
	ID        uint      `gorm:"primaryKey; autoIncrement"`
	Name      string    `gorm:"not null"`
	Token     string    `gorm:"not null; uniqueIndex:ollama_token_token_index"`
	UserId    uint      `gorm:"not null; index:ollama_token_user_id_index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type TokenUsage struct {
	ID              uint      `gorm:"primarykey"`
	UserId          uint      `gorm:"not null; index:token_usage_user_id_index"`
	OllamaModel     string    `gorm:"not null"`
	Time            time.Time `gorm:"autoCreateTime; index:token_usage_time,"`
	PromptEvalCount int
	EvalCount       int
}

func InitModels(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &OllamaToken{}, &TokenUsage{})
	return err
}
