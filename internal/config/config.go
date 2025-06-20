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

	ContainerFilterMode string // "whitelist" or "blacklist"

	// Blacklist (exclude mode)
	ExcludeContainerNames  []string // Container names to ignore
	ExcludeContainerImages []string // Container images to ignore
	ExcludeContainerIDs    []string // Container IDs to ignore
	ExcludeContainerLabels []string // Container labels to ignore

	// Whitelist (include mode)
	IncludeContainerNames  []string // Container names to ignore
	IncludeContainerImages []string // Container images to ignore
	IncludeContainerIDs    []string // Container IDs to ignore
	IncludeContainerLabels []string // Container labels to ignore
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
		ContainerFilterMode:    getEnv("CONTAINER_FILTER_MODE"),
		ExcludeContainerNames:  splitEnv("EXCLUDE_CONTAINER_NAMES"),
		ExcludeContainerImages: splitEnv("EXCLUDE_CONTAINER_IMAGES"),
		ExcludeContainerIDs:    splitEnv("EXCLUDE_CONTAINER_IDS"),
		ExcludeContainerLabels: splitEnv("EXCLUDE_CONTAINER_LABELS"),
		IncludeContainerNames:  splitEnv("INCLUDE_CONTAINER_NAMES"),
		IncludeContainerImages: splitEnv("INCLUDE_CONTAINER_IMAGES"),
		IncludeContainerIDs:    splitEnv("INCLUDE_CONTAINER_IDS"),
		IncludeContainerLabels: splitEnv("INCLUDE_CONTAINER_LABELS"),
		Fiber: Fiber{
			Port: getEnvAsInt("SERVER_PORT"),
		},
	}
}
