package envfile

import (
	"fmt"
	"strings"
)

// FormatExistsResults returns a human-readable list of existence check results.
func FormatExistsResults(results []ExistsResult, colorize bool) string {
	if len(results) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(formatExistsLine(r, colorize))
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatExistsLine(r ExistsResult, colorize bool) string {
	if r.Exists {
		val := r.Value
		if r.Masked {
			val = "***"
		}
		line := fmt.Sprintf("  ✔ %s=%s", r.Key, val)
		if colorize {
			return "\033[32m" + line + "\033[0m"
		}
		return line
	}
	line := fmt.Sprintf("  ✘ %s (not found)", r.Key)
	if colorize {
		return "\033[31m" + line + "\033[0m"
	}
	return line
}

// FormatExistsSummary returns a one-line summary of found/missing counts.
func FormatExistsSummary(s ExistsSummary) string {
	return fmt.Sprintf("%d found, %d missing", s.Found, s.Missing)
}
