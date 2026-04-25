package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatDefaultResults returns a human-readable summary of SetDefault results.
func FormatDefaultResults(results []DefaultResult, colorize bool) string {
	if len(results) == 0 {
		return "(no defaults processed)\n"
	}
	sorted := make([]DefaultResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })

	var sb strings.Builder
	for _, r := range sorted {
		sb.WriteString(formatDefaultLine(r, colorize))
	}
	return sb.String()
}

func formatDefaultLine(r DefaultResult, colorize bool) string {
	displayVal := r.Value
	if IsSensitive(r.Key) {
		displayVal = "[redacted]"
	}
	var prefix, reset string
	if colorize {
		reset = "\033[0m"
		switch r.Status {
		case "applied":
			prefix = "\033[32m" // green
		case "skipped":
			prefix = "\033[33m" // yellow
		case "unchanged":
			prefix = "\033[90m" // grey
		}
	}
	switch r.Status {
	case "applied":
		return fmt.Sprintf("%s  + %s = %s%s\n", prefix, r.Key, displayVal, reset)
	case "skipped":
		return fmt.Sprintf("%s  ~ %s (skipped, has value)%s\n", prefix, r.Key, reset)
	case "unchanged":
		return fmt.Sprintf("%s  = %s (unchanged)%s\n", prefix, r.Key, reset)
	default:
		return fmt.Sprintf("  ? %s\n", r.Key)
	}
}

// FormatDefaultSummary returns a one-line summary string.
func FormatDefaultSummary(sum DefaultSummary) string {
	return fmt.Sprintf("defaults: %d applied, %d skipped, %d unchanged\n",
		sum.Applied, sum.Skipped, sum.Unchanged)
}
