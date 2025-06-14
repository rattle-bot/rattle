package telegram

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/ilyxenc/rattle/internal/config"
	"github.com/ilyxenc/rattle/internal/logger"
)

var client *resty.Client

var (
	baseURL string // Base URL for Telegram Bot API
	chatID  string // Target chat ID for sending notifications
)

// Init initializes the Telegram client and configures retry behavior
func Init() {
	client = resty.New().
		SetRetryCount(5).                     // Number of retry attempts
		SetRetryWaitTime(1 * time.Second).    // Minimum wait between retries
		SetRetryMaxWaitTime(5 * time.Second). // Maximum wait time between retries
		AddRetryCondition(func(r *resty.Response, err error) bool {
			// Retry on network errors or 5xx HTTP status codes
			if err != nil {
				return true
			}
			return r.StatusCode() >= 500
		})

	baseURL = fmt.Sprintf("https://api.telegram.org/bot%s", config.Cfg.BotToken)
	chatID = config.Cfg.ChatID

	logger.Log.Debugf("Telegram initialized for chat %s", chatID)
}

// SendPlainText sends a MarkdownV2-formatted text message to the configured Telegram chat
func SendPlainText(msg string) {
	msg = cleanUTF8(msg) // Sanitize message to ensure it's valid UTF-8

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"chat_id":    chatID,
			"text":       msg,
			"parse_mode": "MarkdownV2", // Enables MarkdownV2 formatting
		}).
		SetHeader("Content-Type", "application/json").
		Get(baseURL + "/sendMessage")

	if err != nil {
		logger.Log.Errorf("Failed to send Telegram message: %v", err)
		return
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Log.Errorf("Telegram responded with status %d: %s", resp.StatusCode(), resp.String())
	}
}
