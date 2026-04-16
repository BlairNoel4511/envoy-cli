package envfile

import (
	"fmt"
	"strings"
)

// FormatGetResults returns a formatted string for a slice of GetResults.
func FormatGetResults(results []GetResult, colorize bool) string {
	if len(results) == 0 {
		return "(no keys requested)"
	}
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(formatGetLine(r, colorize))
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

func formatGetLine(r GetResult, colorize bool) string {
	if !r.Found {
		if colorize {
			return fmt.Sprintf("\033[33m- %s\033[0m (not found)", r.Key)
		}
		return fmt.Sprintf("- %s (not found)", r.Key)
	}
	val := r.Value
	if r.Redacted {
		if colorize {
			return fmt.Sprintf("\033[36m%s\033[0m=\033[31m%s\033[0m", r.Key, val)
		}
		return fmt.Sprintf("%s=%s", r.Key, val)
	}
	if colorize {
		return fmt.Sprintf("\033[36m%s\033[0m=%s", r.Key, val)
	}
	return fmt.Sprintf("%s=%s", r.Key, val)
}
