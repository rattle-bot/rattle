package scanner

import (
	"strings"

	"github.com/ilyxenc/rattle/internal/config"
	"github.com/ilyxenc/rattle/internal/docker"
)

// cleanLine trims extra spaces, invisible characters, and line breaks
func cleanLine(line string) string {
	line = strings.TrimRight(line, " \t\r\n\u00A0\u200B\u202F") // Trailing junk
	return strings.TrimLeft(line, " \t\u00A0\u200B\u202F")      // Leading junk
}

// shouldIgnoreContainer returns true if the container should be excluded from scanning based on its name, image, or ID.
// The check is case-insensitive and supports partial matches for name and image, and prefix match for ID
func shouldIgnoreContainer(ci docker.ContainerInfo) bool {
	name := strings.ToLower(ci.Name)
	image := strings.ToLower(ci.Image)
	id := strings.ToLower(ci.ID)

	// Check if container name matches any exclusion pattern
	for _, n := range config.Cfg.ExcludeContainerNames {
		if strings.Contains(name, strings.ToLower(n)) {
			return true
		}
	}

	// Check if container image matches any exclusion pattern
	for _, img := range config.Cfg.ExcludeContainerImages {
		if strings.Contains(image, strings.ToLower(img)) {
			return true
		}
	}

	// Check if container ID starts with any excluded ID prefix
	for _, i := range config.Cfg.ExcludeContainerIDs {
		if strings.HasPrefix(id, strings.ToLower(i)) {
			return true
		}
	}
	
	return false
}
