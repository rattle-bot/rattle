package scanner

import "strings"

// cleanLine trims extra spaces, invisible characters, and line breaks
func cleanLine(line string) string {
	line = strings.TrimRight(line, " \t\r\n\u00A0\u200B\u202F") // Trailing junk
	return strings.TrimLeft(line, " \t\u00A0\u200B\u202F")      // Leading junk
}
