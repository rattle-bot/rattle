package database

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ilyxenc/rattle/internal/config"
	"github.com/ilyxenc/rattle/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB gorm connector
var DB *gorm.DB

func Connect(path string) error {
	// Ensure the directory for the database file exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	var logLevel logger.LogLevel
	if config.Cfg.Env == "local" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	return err
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&User{}, &models.LogExclusion{}, &models.Chat{}, &models.ContainerExclusion{},
	)
}

func Initialize() error {
	// Init chat ids to send notifications
	for _, chatID := range config.Cfg.ChatIDs {
		if strings.TrimSpace(chatID) == "" {
			continue
		}

		if err := DB.FirstOrCreate(&models.Chat{}, models.Chat{ChatID: chatID}).Error; err != nil {
			return err
		}
	}

	// Init default, include and exclude patterns for logs
	defaultErrorPatterns := []string{
		`(?i)\berror\b`,
		`(?i)\bpanic\b`,
		`(?i)\bfailed\b`,
		`(?i)\bexception\b`,
		`(?i)\btraceback\b`,
		`(?i)\bunhandledpromiserejection\b`,
		`(?i)\bsegmentation fault\b`,
	}

	for _, pattern := range defaultErrorPatterns {
		if strings.TrimSpace(pattern) == "" {
			continue
		}

		entry := models.LogExclusion{
			Pattern:   pattern,
			MatchType: models.MatchTypeInclude,
			EventType: models.EventTypeError,
		}
		if err := DB.FirstOrCreate(&models.LogExclusion{}, entry).Error; err != nil {
			return err
		}
	}

	for eventType, patterns := range config.Cfg.IncludePatterns {
		for _, pattern := range patterns {
			if strings.TrimSpace(pattern) == "" {
				continue
			}

			entry := models.LogExclusion{
				Pattern:   pattern,
				MatchType: models.MatchTypeInclude,
				EventType: eventType,
			}
			if err := DB.FirstOrCreate(&models.LogExclusion{}, entry).Error; err != nil {
				return err
			}
		}
	}

	for _, pattern := range config.Cfg.ExcludePatterns {
		if strings.TrimSpace(pattern) == "" {
			continue
		}

		entry := models.LogExclusion{
			Pattern:   pattern,
			MatchType: models.MatchTypeExclude,
			EventType: "",
		}
		if err := DB.FirstOrCreate(&models.LogExclusion{}, entry).Error; err != nil {
			return err
		}
	}

	// Init exclude containers values
	for _, val := range config.Cfg.ExcludeContainerNames {
		if strings.TrimSpace(val) == "" {
			continue
		}

		entry := models.ContainerExclusion{
			Type:  models.ContainerExclusionName,
			Value: val,
		}
		if err := DB.FirstOrCreate(&models.ContainerExclusion{}, entry).Error; err != nil {
			return err
		}
	}

	for _, val := range config.Cfg.ExcludeContainerImages {
		if strings.TrimSpace(val) == "" {
			continue
		}

		entry := models.ContainerExclusion{
			Type:  models.ContainerExclusionImage,
			Value: val,
		}
		if err := DB.FirstOrCreate(&models.ContainerExclusion{}, entry).Error; err != nil {
			return err
		}
	}

	for _, val := range config.Cfg.ExcludeContainerIDs {
		if strings.TrimSpace(val) == "" {
			continue
		}

		entry := models.ContainerExclusion{
			Type:  models.ContainerExclusionID,
			Value: val,
		}
		if err := DB.FirstOrCreate(&models.ContainerExclusion{}, entry).Error; err != nil {
			return err
		}
	}

	return nil
}
