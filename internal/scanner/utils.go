package scanner

import (
	"strings"

	"github.com/ilyxenc/rattle/internal/docker"
	"github.com/ilyxenc/rattle/internal/managers"
	"github.com/ilyxenc/rattle/internal/models"
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
	if matchesAny(name, managers.Containers.All(models.ContainerExclusionName), strings.Contains) {
		return true
	}

	// Check if container image matches any exclusion pattern
	if matchesAny(image, managers.Containers.All(models.ContainerExclusionImage), strings.Contains) {
		return true
	}

	// Check if container ID starts with any excluded ID prefix
	if matchesAny(id, managers.Containers.All(models.ContainerExclusionID), strings.HasPrefix) {
		return true
	}

	// Check if container labels matches any exclusion pattern
	for key, val := range ci.Labels {
		label := strings.ToLower(key + "=" + val)

		if matchesAny(label, managers.Containers.All(models.ContainerExclusionLabel), strings.Contains) {
			return true
		}
	}

	return false
}

// matchesAny checks if the value matches any pattern using the given matcher function
func matchesAny(value string, patterns []string, matchFn func(string, string) bool) bool {
	for _, p := range patterns {
		if matchFn(value, p) {
			return true
		}
	}
	return false
}
