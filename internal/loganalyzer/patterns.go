package loganalyzer

import (
	"strings"

	"github.com/ilyxenc/rattle/internal/managers"
)

// IsLogError returns true if the provided log line matches any known error pattern
func IsLogError(line string) bool {
	line = strings.TrimSpace(line)

	// Check exclusion patterns first
	for _, re := range managers.Logs.Exclude() {
		if re.MatchString(line) {
			return false
		}
	}

	// Check inclusion patterns for "error" type
	for _, re := range managers.Logs.Include("error") {
		if re.MatchString(line) {
			return true
		}
	}

	return false
}

// DetectEventType returns the matching event type for the given line, or empty string if it matches nothing or is excluded
func DetectEventType(line string) string {
	if line == "" {
		return ""
	}

	// First check if line is excluded
	for _, re := range managers.Logs.Exclude() {
		if re.MatchString(line) {
			return ""
		}
	}

	// Now check which event type matches
	for _, eventType := range managers.Logs.KnownEventTypes() {
		for _, re := range managers.Logs.Include(eventType) {
			if re.MatchString(line) {
				return eventType
			}
		}
	}

	// No match
	return ""
}
