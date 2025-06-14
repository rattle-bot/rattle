package docker

import (
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/ilyxenc/rattle/internal/logger"
)

// ContainerInfo holds simplified metadata for a Docker container
type ContainerInfo struct {
	ID      string // Full container ID
	Name    string // Container name without leading slash
	Image   string // Image name (e.g. redis:latest)
	ImageID string // Full image ID
	ShortID string // First 12 characters of the container ID
}

// NewContainerInfo safely extracts basic info from container summary
func NewContainerInfo(c container.Summary) ContainerInfo {
	name := ""
	if len(c.Names) > 0 {
		name = strings.TrimPrefix(c.Names[0], "/") // Remove leading slash
	}

	return ContainerInfo{
		ID:      c.ID,
		Name:    name,
		Image:   c.Image,
		ImageID: c.ImageID,
		ShortID: shortID(c.ID),
	}
}

// shortID returns the first 12 characters of a container ID
func shortID(id string) string {
	if len(id) >= 12 {
		return id[:12]
	}
	return id
}

// NewClient creates Docker client using environment and handles version negotiation
func NewClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Log.Panicf("Failed to create docker client: %v", err)
	}

	return cli
}
