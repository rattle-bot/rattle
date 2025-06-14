package logger

import (
	"log"

	"github.com/ilyxenc/rattle/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is the global sugared logger instance used throughout the application
var Log *zap.SugaredLogger

// parseLogLevel converts a string to a zapcore.Level.
// Falls back to InfoLevel if the string is unrecognized
func parseLogLevel(s string) zapcore.Level {
	switch s {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Init initializes the global Zap logger based on config
func Init() {
	level := parseLogLevel(config.Cfg.LogLevel)

	var cfg zap.Config
	if config.Cfg.Env == "local" || config.Cfg.Env == "dev" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		log.Panicf("Failed to init zap: %v", err)
	}

	Log = logger.Sugar()
}
