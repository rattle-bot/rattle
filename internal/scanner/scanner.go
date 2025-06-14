package scanner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/ilyxenc/rattle/internal/logger"
)

// Start begins streaming logs from the container.
// It automatically reconnects on failure unless MaxRetry is reached
func (s *LogScanner) Start(ctx context.Context) error {
	since := s.Since
	retries := 0

	for {
		// Attempt a single log stream session
		err := s.streamOnce(ctx, since)
		if err == nil || ctx.Err() != nil {
			return err // Exit on success or if context was cancelled
		}

		logger.Log.Warnw("log stream error — retrying",
			"container", s.Container.Name,
			"error", err,
			"retry", retries+1,
		)

		retries++
		if s.MaxRetry > 0 && retries >= s.MaxRetry {
			return fmt.Errorf("max retry reached for %s: %w", s.Container.Name, err)
		}

		time.Sleep(s.ReconnectDelay)
		since = time.Now() // Update `since` to avoid getting duplicates
	}
}

// streamOnce connects to the container logs and reads them line by line
func (s *LogScanner) streamOnce(ctx context.Context, since time.Time) error {
	reader, err := s.Client.ContainerLogs(ctx, s.Container.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Since:      fmt.Sprintf("%d", since.Unix()), // Start from given timestamp
		Timestamps: false,
	})
	if err != nil {
		return err
	}
	defer reader.Close()

	// Pipe to capture clean stdout/stderr from stdcopy
	pr, pw := io.Pipe()

	// Strip Docker's multiplexed headers and write to pipe
	go func() {
		defer pw.Close()
		_, _ = stdcopy.StdCopy(pw, pw, reader)
	}()

	scanner := bufio.NewScanner(pr)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil // Exit gracefully if cancelled
		default:
			line := cleanLine(scanner.Text())
			if line != "" && s.OnLog != nil {
				s.OnLog(s.Container, line) // Forward log line to callback
			}
		}
	}

	// Return any unexpected scanner error (excluding EOF)
	if err := scanner.Err(); err != nil && err != io.EOF {
		return err
	}

	return nil
}
