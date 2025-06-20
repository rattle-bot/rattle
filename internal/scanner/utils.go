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

// shouldIgnoreContainer determines whether a container should be excluded from scanning, based on the current filtering mode (whitelist or blacklist) and matching rules
//
// Matching logic:
//   - In "whitelist" mode: the container is ignored unless it matches at least one entry
//     from the include list (ID prefix, label, name or image)
//   - In "blacklist" mode: the container is ignored if it matches any entry
//     from the exclude list (ID prefix, label, name or image)
func shouldIgnoreContainer(ci docker.ContainerInfo) bool {
	name := strings.ToLower(ci.Name)
	image := strings.ToLower(ci.Image)
	id := strings.ToLower(ci.ID)
	labels := make([]string, 0, len(ci.Labels))
	for key, val := range ci.Labels {
		labels = append(labels, strings.ToLower(key+"="+val))
	}

	mode := managers.Mode.Get()

	// Whitelist mode: must match at least one
	if mode == models.Whitelist {
		if matchesAny(id, managers.Containers.All(models.ContainerID, mode), strings.Contains) {
			return false
		}

		if anyLabelMatches(labels, managers.Containers.All(models.ContainerLabel, mode)) {
			return false
		}

		if matchesAny(name, managers.Containers.All(models.ContainerName, mode), strings.Contains) {
			return false
		}

		if matchesAny(image, managers.Containers.All(models.ContainerImage, mode), strings.Contains) {
			return false
		}

		return true // Not in whitelist → ignore
	}

	// Blacklist mode: must NOT match any
	if matchesAny(id, managers.Containers.All(models.ContainerID, mode), strings.Contains) {
		return true
	}

	if anyLabelMatches(labels, managers.Containers.All(models.ContainerLabel, mode)) {
		return true
	}

	if matchesAny(name, managers.Containers.All(models.ContainerName, mode), strings.Contains) {
		return true
	}

	if matchesAny(image, managers.Containers.All(models.ContainerImage, mode), strings.Contains) {
		return true
	}

	return false // Not in blacklist → allow
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

// anyLabelMatches checks if labels values matches any pattern
func anyLabelMatches(actualLabels []string, filters []string) bool {
	for _, actual := range actualLabels {
		if matchesAny(actual, filters, strings.Contains) {
			return true
		}
	}
	return false
}
