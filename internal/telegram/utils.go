package telegram

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ilyxenc/rattle/internal/docker"
)

// formatMessage formats an error message for Telegram using MarkdownV2 code block
func formatMessage(eventType, errText string) string {
	cleaned := cleanUTF8(errText)
	return fmt.Sprintf("\n\n```%s\n%s\n```", escapeMarkdownV2(eventType), escapeMarkdownV2(cleaned))
}

// formatMeta returns formatted container metadata with timestamp, used as part of notifications
func formatMeta(ci docker.ContainerInfo) string {
	return fmt.Sprintf(
		"\n\n📦 ID: `%s`\nName: `%s`\nImage: `%s`\n\n|| %s ||",
		ci.ShortID, ci.Name, ci.Image, escapeMarkdownV2(time.Now().Format("2006-01-02T15:04:05.000Z07:00")), // With milliseconds
	)
}

// formatContainersSummary returns formatted information about active containers
func formatContainersSummary(containers []docker.ContainerInfo) string {
	if len(containers) == 0 {
		return "📦 No active containers running"
	}

	msg := fmt.Sprintf("📊 *%d active containers:*\n\n", len(containers))
	for _, ci := range containers {
		msg += fmt.Sprintf("\\- `%s`: %s\n", ci.ShortID, escapeMarkdownV2(ci.Name))
	}
	return msg
}

// cleanUTF8 removes invalid UTF-8 runes from the input string to ensure Telegram accepts the message
func cleanUTF8(input string) string {
	if utf8.ValidString(input) {
		return input
	}

	// Filter only valid runes
	out := make([]rune, 0, len(input))
	for i, r := range input {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(input[i:])
			if size == 1 {
				continue // Skip invalid rune
			}
		}
		out = append(out, r)
	}
	return string(out)
}

// escapeMarkdownV2 escapes all special MarkdownV2 characters in a string to prevent formatting issues or Telegram API errors
func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		`_`, `\_`,
		`*`, `\*`,
		`[`, `\[`,
		`]`, `\]`,
		`(`, `\(`,
		`)`, `\)`,
		`~`, `\~`,
		"`", "\\`",
		`>`, `\>`,
		`#`, `\#`,
		`+`, `\+`,
		`-`, `\-`,
		`=`, `\=`,
		`|`, `\|`,
		`{`, `\{`,
		`}`, `\}`,
		`.`, `\.`,
		`!`, `\!`,
	)
	return replacer.Replace(text)
}

// EventEmoji returns emoji for title based on event type
func EventEmoji(eventType string) string {
	switch eventType {
	case "error":
		return "❌"
	case "warning":
		return "⚠️"
	case "success":
		return "✅"
	case "info":
		return "ℹ️"
	case "critical":
		return "🚨"
	default:
		return "📦"
	}
}

func FormatEventTitle(eventType, containerName string) string {
	switch eventType {
	case "error":
		return fmt.Sprintf("❌ *Error in container:* `%s`", containerName)
	case "warning":
		return fmt.Sprintf("⚠️ *Warning in container:* `%s`", containerName)
	case "success":
		return fmt.Sprintf("✅ *Success in container:* `%s`", containerName)
	case "info":
		return fmt.Sprintf("ℹ️ *Info from container:* `%s`", containerName)
	case "critical":
		return fmt.Sprintf("🚨 *Critical event in container:* `%s`", containerName)
	default:
		return fmt.Sprintf("📦 *Log from container:* `%s`", containerName)
	}
}
