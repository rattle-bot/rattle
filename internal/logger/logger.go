package logger

import (
	"os"

	"github.com/ilyxenc/rattle/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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
func Init(path string) {
	level := parseLogLevel(config.Cfg.LogLevel)

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path, // Relative path from project root (main.go)
		MaxSize:    10,   // Megabytes
		MaxBackups: 3,    // Number of old log files to keep
		MaxAge:     28,   // Max age in days before deletion
		Compress:   true, // Compress old log files
	})

	// Encoder configuration
	cfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Color for console
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Encoders for human-readable (console) and JSON (file) outputs
	consoleEncoder := zapcore.NewConsoleEncoder(cfg)
	fileEncoder := zapcore.NewJSONEncoder(cfg)

	var cores []zapcore.Core

	// In dev/local, also log to console
	if config.Cfg.Env == "local" || config.Cfg.Env == "dev" {
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level))
	}

	// Always log to file
	cores = append(cores, zapcore.NewCore(fileEncoder, fileWriter, level))

	// Combine all logging targets
	core := zapcore.NewTee(cores...)

	// Build logger with caller info and stacktrace for warnings+
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Assign to global logger
	Log = logger.Sugar()
}
