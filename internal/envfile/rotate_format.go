package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatRotateSummary returns a human-readable summary of a rotation run.
// When colorize is true ANSI escape codes are added.
func FormatRotateSummary(summary RotateSummary, colorize bool) string {
	var sb strings.Builder

	green := func(s string) string {
		if colorize {
			return "\033[32m" + s + "\033[0m"
		}
		return s
	}
	yellow := func(s string) string {
		if colorize {
			return "\033[33m" + s + "\033[0m"
		}
		return s
	}

	// Sort for deterministic output.
	sorted := make([]RotateResult, len(summary.Results))
	copy(sorted, summary.Results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, r := range sorted {
		if r.Skipped {
			sb.WriteString(yellow(fmt.Sprintf("  ~ %s (skipped: %s)\n", r.Key, r.Reason)))
		} else {
			maskedOld := maskRotateValue(r.Key, r.OldValue)
			maskedNew := maskRotateValue(r.Key, r.NewValue)
			sb.WriteString(green(fmt.Sprintf("  ↻ %s: %s → %s\n", r.Key, maskedOld, maskedNew)))
		}
	}

	sb.WriteString(fmt.Sprintf("\nRotated: %d  Skipped: %d\n", summary.Rotated, summary.Skipped))
	return sb.String()
}

func maskRotateValue(key, value string) string {
	if IsSensitive(key) {
		return "[redacted]"
	}
	return value
}
