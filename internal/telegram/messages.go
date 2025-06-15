package telegram

import (
	"fmt"

	"github.com/ilyxenc/rattle/internal/config"
	"github.com/ilyxenc/rattle/internal/docker"
)

// NotificationType defines the type of event being reported to Telegram
type NotificationType string

const (
	NotificationContainerStart         NotificationType = "container_start"           // Sent when a container starts
	NotificationContainerStop          NotificationType = "container_stop"            // Sent when a container stops normally
	NotificationLogEvent               NotificationType = "log_event"                 // Sent when an event is detected in logs
	NotificationContainerStopWithError NotificationType = "container_stop_with_error" // Sent when a container stops unexpectedly with error
	NotificationShutDownRattle         NotificationType = "shut_down_rattle"          // Sent when Rattle is shutting down
	NotificationStartedRattle          NotificationType = "started_rattle"            // Sent when Rattle starts
	NotificationContainersSummary      NotificationType = "containers_summary"        // Sent when Rattle starts and find containers
)

// Notification represents the structure of a message to be sent to Telegram
type Notification struct {
	Type       NotificationType       // The type of event
	EventType  string                 // error, info, success, warning, critical
	Title      string                 // Notification template title
	Details    string                 // Optional details (e.g., error message or log content)
	Container  docker.ContainerInfo   // Metadata about the container related to the event
	Containers []docker.ContainerInfo // For summary events like containers list
}

// Notify sends a formatted notification to the configured Telegram chat
func Notify(n Notification) {
	msg := RenderNotification(n)
	SendPlainText(msg)
}

// RenderNotification formats a notification message based on its type
func RenderNotification(n Notification) string {
	c := n.Container

	switch n.Type {
	case NotificationLogEvent:
		title := FormatEventTitle(n.EventType, escapeMarkdownV2(n.Container.Name))
		return title + formatMessage(n.EventType, n.Details) + formatMeta(c)
	case NotificationContainerStart:
		return fmt.Sprintf("✅ *Container started:* `%s`", c.Name) + formatMeta(c)
	case NotificationContainerStop:
		return fmt.Sprintf("🛑 *Container stopped:* `%s`", c.Name) + formatMeta(c)
	case NotificationContainerStopWithError:
		return fmt.Sprintf("🛑 *Container stopped with error:* `%s`", c.Name) + formatMeta(c)
	case NotificationShutDownRattle:
		return fmt.Sprintf("🛑 *Rattle is shutting down%s*", escapeMarkdownV2("..."))
	case NotificationStartedRattle:
		return fmt.Sprintf("🚀 Rattle started in *%s* mode", config.Cfg.Env)
	case NotificationContainersSummary:
		return formatContainersSummary(n.Containers)
	default:
		return "📦 Unknown notification type"
	}
}
