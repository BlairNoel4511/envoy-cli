package envfile

import "strings"

// RedactedValue is the placeholder used in place of sensitive values.
const RedactedValue = "***REDACTED***"

// sensitiveKeyPatterns lists substrings that indicate a sensitive key.
var sensitiveKeyPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"AUTH",
	"CREDENTIAL",
	"DSN",
}

// IsSensitive reports whether a key is considered sensitive.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range sensitiveKeyPatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// Redact returns a copy of the EnvFile with sensitive values replaced.
func (e *EnvFile) Redact() *EnvFile {
	redacted := &EnvFile{
		Path:    e.Path,
		Entries: make([]Entry, len(e.Entries)),
	}
	for i, entry := range e.Entries {
		if entry.Key != "" && IsSensitive(entry.Key) {
			redacted.Entries[i] = Entry{
				Key:   entry.Key,
				Value: RedactedValue,
			}
		} else {
			redacted.Entries[i] = entry
		}
	}
	return redacted
}
