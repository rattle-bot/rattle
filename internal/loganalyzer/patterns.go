package loganalyzer

import (
	"regexp"
	"sync"
)

var (
	// errorPatterns is a list of regular expressions used to detect error messages in log lines
	errorPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)\berror\b`),
		regexp.MustCompile(`(?i)\bpanic\b`),
		regexp.MustCompile(`(?i)\bfailed\b`),
		regexp.MustCompile(`(?i)\bexception\b`),
		regexp.MustCompile(`(?i)\btraceback\b`),
		regexp.MustCompile(`(?i)\bunhandledpromiserejection\b`),
		regexp.MustCompile(`(?i)\bsegmentation fault\b`),
		regexp.MustCompile(`(?i)(^|\s|:|]|\[)ошибка(:|\s|$)`), // Russian: "ошибка"
	}
	mu sync.RWMutex // Protects access to errorPatterns
)

// IsLogError returns true if the provided log line matches any known error pattern
func IsLogError(line string) bool {
	mu.RLock() // Read only
	defer mu.RUnlock()

	if len(errorPatterns) == 0 {
		return false
	}

	for _, re := range errorPatterns {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

// AddErrorPattern allows dynamically adding a new error pattern (e.g., user-defined error string)
func AddErrorPattern(pattern string) error {
	mu.Lock()
	defer mu.Unlock()

	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	errorPatterns = append(errorPatterns, re)
	return nil
}

// ClearPatterns resets the list of error patterns. Useful for testing or reconfiguration
func ClearPatterns() {
	mu.Lock()
	defer mu.Unlock()

	errorPatterns = []*regexp.Regexp{}
}
