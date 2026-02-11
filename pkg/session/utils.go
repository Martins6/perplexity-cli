package session

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GetSessionsDir returns the directory where sessions are stored
func GetSessionsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if can't get home
		return ".pplx/sessions"
	}
	return filepath.Join(home, ".pplx", "sessions")
}

// EnsureSessionsDir creates the sessions directory if it doesn't exist
func EnsureSessionsDir() error {
	sessionsDir := GetSessionsDir()
	return os.MkdirAll(sessionsDir, 0755)
}

// generateSessionID generates a unique session ID based on timestamp
// Format: 20060102-150405.000 (date-time.milliseconds)
func generateSessionID(t time.Time) string {
	return t.Format("20060102-150405.000")
}

// GenerateSessionFilename generates a full filename for a session
func GenerateSessionFilename(id string) string {
	return filepath.Join(GetSessionsDir(), id+".json")
}

// ParseSessionID parses a session ID from a filename
func ParseSessionID(filename string) string {
	base := filepath.Base(filename)
	// Remove .json extension
	if len(base) > 5 && base[len(base)-5:] == ".json" {
		return base[:len(base)-5]
	}
	return base
}

// IsValidSessionFile checks if a file is a valid session file
func IsValidSessionFile(filename string) bool {
	base := filepath.Base(filename)
	// Check if it ends with .json
	if len(base) < 6 || base[len(base)-5:] != ".json" {
		return false
	}
	// Just check that it's a .json file - be lenient with the ID format
	return true
}

// TruncateQuery truncates a query for display purposes
func TruncateQuery(query string, maxLen int) string {
	if len(query) <= maxLen {
		return query
	}
	if maxLen <= 3 {
		return query[:maxLen]
	}
	return query[:maxLen-3] + "..."
}

// FormatSessionTime formats a session time for display
func FormatSessionTime(t time.Time) string {
	return t.Format("Jan 02, 2006 15:04:05")
}

// Debugf prints debug information if DEBUG environment variable is set
func Debugf(format string, args ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
	}
}

// base62Chars is the character set used for Base62 encoding
const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// encodeBase62 encodes an int64 as a Base62 string
func encodeBase62(n int64) string {
	if n == 0 {
		return string(base62Chars[0])
	}

	var result []byte
	base := int64(len(base62Chars))
	for n > 0 {
		result = append(result, base62Chars[n%base])
		n /= base
	}

	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// GenerateShortID generates a short, unique session ID using Base62-encoded Unix timestamp
// The ID is based on Unix timestamp in milliseconds, ensuring uniqueness and monotonic ordering
func GenerateShortID(t time.Time) string {
	return encodeBase62(t.UnixMilli())
}
