package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Postgres struct {
	Port     int
	Host     string
	User     string
	Password string
	Database string
}

type Fiber struct {
	Port int
}

// Config holds all environment-based configuration for the application
type Config struct {
	BotToken string   // Telegram bot token
	ChatIDs  []string // Telegram chat IDs to send messages to
	LogLevel string   // Log level: debug, info, warn, error
	Env      string   // Application environment: local, dev, prod, etc
	Postgres Postgres
	Fiber    Fiber

	IncludePatterns map[string][]string // Key = eventType
	ExcludePatterns []string            // Regex patterns to exclude from log detection

	ExcludeContainerNames  []string // Container names to ignore
	ExcludeContainerImages []string // Container images to ignore
	ExcludeContainerIDs    []string // Container IDs to ignore
}

// Cfg is the global config instance accessible throughout the app
var Cfg *Config

// Load loads environment variables into the global config struct
func Load() {
	// Try to load from .env file, fallback to system env if not found
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize the global config from required env vars
	Cfg = &Config{
		BotToken: getEnv("TELEGRAM_BOT_TOKEN"),
		ChatIDs:  splitEnv("TELEGRAM_CHAT_IDS"),
		LogLevel: getEnv("LOG_LEVEL"),
		Env:      getEnv("APP_ENV"),
		Postgres: Postgres{
			Port:     getEnvAsInt("POSTGRES_PORT"),
			Host:     getEnv("POSTGRES_HOST"),
			User:     getEnv("POSTGRES_USER"),
			Password: getEnv("POSTGRES_PASSWORD"),
			Database: getEnv("POSTGRES_DB"),
		},
		IncludePatterns: map[string][]string{
			"error":    splitEnv("INCLUDE_PATTERNS_ERROR"),
			"success":  splitEnv("INCLUDE_PATTERNS_SUCCESS"),
			"info":     splitEnv("INCLUDE_PATTERNS_INFO"),
			"warning":  splitEnv("INCLUDE_PATTERNS_WARNING"),
			"critical": splitEnv("INCLUDE_PATTERNS_CRITICAL"),
		},
		ExcludePatterns:        splitEnv("EXCLUDE_PATTERNS"),
		ExcludeContainerNames:  splitEnv("EXCLUDE_CONTAINER_NAMES"),
		ExcludeContainerImages: splitEnv("EXCLUDE_CONTAINER_IMAGES"),
		ExcludeContainerIDs:    splitEnv("EXCLUDE_CONTAINER_IDS"),
		Fiber: Fiber{
			Port: getEnvAsInt("FIBER_PORT"),
		},
	}
}
