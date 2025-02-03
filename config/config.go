package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

var ServerAddr string

var DatabasePath string

var LoggingLevel slog.Level
var LoggingOutput string

var JwtSecret []byte

var OllamaHost string
var OllamaTimeout int

func ReadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	initValue()
}

func initValue() {
	ServerAddr = GetStringWithDefault("server.address", ":8080")

	DatabasePath = GetStringWithDefault("database.url", "safe_ollama.db")

	loggingLevel := strings.ToLower(GetStringWithDefault("server.logging.level", "info"))
	switch loggingLevel {
	case "debug":
		LoggingLevel = slog.LevelDebug
	case "info":
		LoggingLevel = slog.LevelInfo
	case "warn":
		LoggingLevel = slog.LevelWarn
	case "error":
		LoggingLevel = slog.LevelError
	default:
		slog.Warn("Invalid logging level, using default value \"info\"")
		LoggingLevel = slog.LevelInfo
	}
	LoggingOutput = GetStringWithDefault("server.logging.output", "stdout")

	JwtSecret = []byte(GetStringWithDefault("server.jwtkey", "jwt_secret_key"))

	OllamaHost = GetStringWithDefault("ollama.url", "http://localhost:11434")
	OllamaTimeout = GetIntWithDefault("ollama.timeout", 300)
}

func GetStringWithDefault(key string, defaultValue string) string {
	if value := viper.GetString(key); viper.IsSet(key) && value != "" {
		return value
	} else {
		slog.Warn(fmt.Sprintf("Config \"%s\" not found, using default value \"%s\"", key, defaultValue))
		return defaultValue
	}
}

func GetIntWithDefault(key string, defaultValue int) int {
	if value := viper.GetInt(key); viper.IsSet(key) {
		return value
	} else {
		slog.Warn(fmt.Sprintf("Config \"%s\" not found, using default value \"%d\"", key, defaultValue))
		return defaultValue
	}
}
