package loganalyzer

import (
	"regexp"
	"sync"

	"github.com/ilyxenc/rattle/internal/config"
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
	}
	includePatterns []*regexp.Regexp
	excludePatterns []*regexp.Regexp
	mu              sync.RWMutex // Protects access to errorPatterns
)

// InitCustomPatterns compiles include/exclude patterns from config
func InitCustomPatterns() {
	mu.Lock()
	defer mu.Unlock()

	includePatterns = compileMany(config.Cfg.IncludeError)
	excludePatterns = compileMany(config.Cfg.ExcludeError)
}

// compileMany takes a slice of string patterns and returns a slice of compiled regular expressions
//
// Invalid patterns are silently skipped to avoid breaking the application during runtime.
// This function is typically used to initialize include/exclude error matchers from config
//
// Example:
//	compileMany([]string{`(?i)error`, `timeout`}) -> []*regexp.Regexp{...}
func compileMany(patterns []string) []*regexp.Regexp {
	var result []*regexp.Regexp
	for _, p := range patterns {
		if re, err := regexp.Compile(p); err == nil {
			result = append(result, re)
		}
	}
	return result
}

// IsLogError returns true if the provided log line matches any known error pattern
func IsLogError(line string) bool {
	mu.RLock() // Read only
	defer mu.RUnlock()

	// Include takes priority — if specified, only match those
	if len(includePatterns) > 0 {
		for _, re := range includePatterns {
			if re.MatchString(line) {
				return true
			}
		}
	}

	// If excluded explicitly, skip
	for _, re := range excludePatterns {
		if re.MatchString(line) {
			return false
		}
	}

	// Fallback to default error patterns
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
