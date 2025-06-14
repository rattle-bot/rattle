package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all environment-based configuration for the application
type Config struct {
	BotToken string // Telegram bot token
	ChatID   string // Telegram chat ID to send messages to
	LogLevel string // Log level: debug, info, warn, error
	Env      string // Application environment: local, dev, prod, etc
}

// Cfg is the global config instance accessible throughout the app
var Cfg *Config

// getEnv returns the value of an environment variable or panics if it's not set
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("Required environment variable not set: %v", key)
	}

	return value
}

// Load loads environment variables into the global config struct
func Load() {
	// Try to load from .env file, fallback to system env if not found
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize the global config from required env vars
	Cfg = &Config{
		BotToken: getEnv("TELEGRAM_BOT_TOKEN"),
		ChatID:   getEnv("TELEGRAM_CHAT_ID"),
		LogLevel: getEnv("LOG_LEVEL"),
		Env:      getEnv("APP_ENV"),
	}
}
