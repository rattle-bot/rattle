package config

import (
	"log"
	"os"
	"strings"
)

// getEnv returns the value of an environment variable or panics if it's not set
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("Required environment variable not set: %v", key)
	}

	return value
}

// splitEnv parses comma-separated string into trimmed slice of strings
func splitEnv(key string) []string {
	val := os.Getenv(key)
	if val == "" {
		return []string{}
	}
	parts := strings.Split(val, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
