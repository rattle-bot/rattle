package telegram

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ilyxenc/rattle/internal/docker"
)

// formatErrorMessage formats an error message for Telegram using MarkdownV2 code block
func formatErrorMessage(errText string) string {
	cleaned := cleanUTF8(errText)
	return fmt.Sprintf("\n\n```Error\n%s\n```", escapeMarkdownV2(cleaned))
}

// formatMeta returns formatted container metadata with timestamp, used as part of notifications
func formatMeta(ci docker.ContainerInfo) string {
	return fmt.Sprintf(
		"\n\n📦 ID: `%s`\nName: `%s`\nImage: `%s`\n\n|| %s ||",
		ci.ShortID, ci.Name, ci.Image, escapeMarkdownV2(time.Now().Format("2006-01-02T15:04:05.000Z07:00")), // With milliseconds
	)
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
