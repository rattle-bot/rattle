package loganalyzer

import (
	"github.com/ilyxenc/rattle/internal/docker"
	"github.com/ilyxenc/rattle/internal/telegram"
)

// AnalyzeLogLine checks if the given log line matches known error patterns.
// If it does, a notification is sent via Telegram
func AnalyzeLogLine(c docker.ContainerInfo, line string) {
	if IsLogError(line) {
		telegram.Notify(telegram.Notification{
			Type:      telegram.NotificationError,
			Details:   line,
			Container: c,
		})
	}
}
