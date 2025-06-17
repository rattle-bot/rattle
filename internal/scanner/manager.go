package scanner

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/ilyxenc/rattle/internal/docker"
	"github.com/ilyxenc/rattle/internal/loganalyzer"
	"github.com/ilyxenc/rattle/internal/logger"
	"github.com/ilyxenc/rattle/internal/telegram"
)

// LogScanManager manages log scanners for all running containers
type LogScanManager struct {
	Ctx      context.Context        // Shared context for cancellation
	Client   *client.Client         // Docker client
	Scanners map[string]*LogScanner // Map of active scanners by container ID
	Mu       sync.Mutex             // Mutex to protect Scanners map
	wg       sync.WaitGroup         // Tracks active scanner goroutines
}

// NewLogScanManager creates a new LogScanManager instance
func NewLogScanManager(ctx context.Context, cli *client.Client) *LogScanManager {
	return &LogScanManager{
		Client:   cli,
		Ctx:      ctx,
		Scanners: make(map[string]*LogScanner),
	}
}

// StartAll starts log scanners for all currently running containers.
// Also starts a background goroutine to watch container lifecycle events
func (m *LogScanManager) StartAll() error {
	containers, err := m.Client.ContainerList(m.Ctx, container.ListOptions{
		All: false, // Only running containers
	})
	if err != nil {
		return err
	}

	// Save active containers for notification
	active := make([]docker.ContainerInfo, 0, len(containers))
	for _, c := range containers {
		ci := docker.NewContainerInfo(c)
		if shouldIgnoreContainer(ci) {
			continue
		}
		m.startScanner(c, true)
		active = append(active, ci)
	}

	telegram.Notify(telegram.Notification{
		Type:       telegram.NotificationContainersSummary,
		Containers: active,
	})

	go m.watchContainerEvents()

	return nil
}

// StopAll stops all active scanners and waits for them to finish
func (m *LogScanManager) StopAll() {
	m.Mu.Lock()
	for _, s := range m.Scanners {
		if s.Cancel != nil {
			s.Cancel()
		}
	}
	m.Mu.Unlock()

	m.wg.Wait() // Wait for all scanner goroutines to complete
}

// startScanner creates and starts a log scanner for the given container.
// If a scanner already exists, it will be stopped and replaced
func (m *LogScanManager) startScanner(c container.Summary, suppressNotify bool) {
	info := docker.NewContainerInfo(c)

	m.Mu.Lock()
	// Cancel and remove existing scanner if present
	if old, exists := m.Scanners[info.ID]; exists {
		if old.Cancel != nil {
			old.Cancel()
		}
		delete(m.Scanners, info.ID)
	}

	// Create new context for this scanner
	ctx, cancel := context.WithCancel(m.Ctx)

	s := &LogScanner{
		Client:         m.Client,
		Container:      info,
		OnLog:          loganalyzer.AnalyzeLogLine,
		Since:          time.Now(),
		ReconnectDelay: 5 * time.Second,
		MaxRetry:       0,
		Cancel:         cancel,
	}
	m.Scanners[info.ID] = s
	m.Mu.Unlock()

	// Notify about container start if not already started
	if !suppressNotify {
		telegram.Notify(telegram.Notification{
			Type:      telegram.NotificationContainerStart,
			Container: info,
		})
	}
	logger.Log.Infof("Started scanner for container %s", info.Name)

	m.wg.Add(1) // Register a new scanner goroutine in the WaitGroup
	go func() {
		defer m.wg.Done() // Signal that this scanner goroutine has finished

		// Start log scanner
		err := s.Start(ctx)

		if err != nil {
			telegram.Notify(telegram.Notification{
				Type:      telegram.NotificationContainerStopWithError,
				Container: info,
			})
			logger.Log.Warnw("Scanner stopped", "container", info.Name, "error", err)
		}

		// Remove scanner from the registry
		m.Mu.Lock()
		delete(m.Scanners, info.ID)
		m.Mu.Unlock()

		// Skip notification if container stopped due to app shutdown
		if errors.Is(ctx.Err(), context.Canceled) {
			return
		}

		// Notify only if not in shutdown mode
		telegram.Notify(telegram.Notification{
			Type:      telegram.NotificationContainerStop,
			Container: info,
		})
		logger.Log.Infof("Scanner removed for container %s", info.Name)
	}()
}

// watchContainerEvents listens for Docker container start/stop/restart events and updates scanners accordingly
func (m *LogScanManager) watchContainerEvents() {
	eventFilter := filters.NewArgs()
	eventFilter.Add("type", "container")
	eventFilter.Add("event", "start")
	eventFilter.Add("event", "die")
	eventFilter.Add("event", "destroy")

	eventsCh, errsCh := m.Client.Events(m.Ctx, events.ListOptions{Filters: eventFilter})

	for {
		select {
		case <-m.Ctx.Done():
			return
		case event := <-eventsCh:
			id := event.Actor.ID
			name := event.Actor.Attributes["name"]

			switch event.Action {
			case "start":
				logger.Log.Infof("Container started: %s", name)

				// Find full container info and start scanner
				containers, err := m.Client.ContainerList(m.Ctx, container.ListOptions{
					All: false,
				})
				if err != nil {
					logger.Log.Errorf("Failed to list containers: %v", err)
					continue
				}

				for _, c := range containers {
					if c.ID == id {
						ci := docker.NewContainerInfo(c)
						if shouldIgnoreContainer(ci) {
							logger.Log.Infof("Ignored container %s due to filters", name)
							continue
						}

						m.startScanner(c, false)
						break
					}
				}
			case "die", "destroy":
				logger.Log.Infof("Container stopped/destroyed: %s", name)

				// Cancel and remove scanner
				m.Mu.Lock()
				if s, ok := m.Scanners[id]; ok {
					if s.Cancel != nil {
						s.Cancel()
					}
					delete(m.Scanners, id)
					logger.Log.Infof("Stopped scanner for container %s", name)
				}
				m.Mu.Unlock()
			}
		case err := <-errsCh:
			logger.Log.Errorf("Docker event error: %v", err)

			return
		}
	}
}
