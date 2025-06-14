package scanner

import (
	"context"
	"time"

	"github.com/docker/docker/client"
	"github.com/ilyxenc/rattle/internal/docker"
)

// OnLogFunc is the type for a log line analyzer callback
type OnLogFunc func(container docker.ContainerInfo, line string)

// LogScanner streams logs from a specific Docker container
//
// It connects to the container's stdout and stderr streams starting from a given timestamp (`Since`).
// For each log line, it calls the `OnLog` callback. If the log stream is interrupted,
// it automatically retries connecting with a specified delay and retry limit
type LogScanner struct {
	Client         *client.Client       // Docker client used to access container logs
	Container      docker.ContainerInfo // Container info like id, name etc
	OnLog          OnLogFunc            // Callback function invoked for each log line received
	Since          time.Time            // Since specifies the starting point for reading container logs. It helps avoid processing old logs or duplicating entries after reconnects
	ReconnectDelay time.Duration        // Delay between reconnection attempts when the log stream fails
	MaxRetry       int                  // Number of times to retry connecting to the log stream before giving up. A value of 0 means unlimited retries
	Cancel         context.CancelFunc   // Cancel function to stop log streaming
}
