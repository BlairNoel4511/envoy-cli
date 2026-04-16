package envfile

import (
	"fmt"
	"strings"
)

// FormatDeleteResults formats delete results into a human-readable string.
func FormatDeleteResults(results []DeleteResult, colorize bool) string {
	if len(results) == 0 {
		return "(no keys specified)"
	}

	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(formatDeleteLine(r, colorize))
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

func formatDeleteLine(r DeleteResult, colorize bool) string {
	if r.Deleted {
		if colorize {
			return fmt.Sprintf("\033[31m- %s\033[0m", r.Key)
		}
		return fmt.Sprintf("- %s", r.Key)
	}
	if colorize {
		return fmt.Sprintf("\033[33m~ %s (%s)\033[0m", r.Key, r.Reason)
	}
	return fmt.Sprintf("~ %s (%s)", r.Key, r.Reason)
}

// FormatDeleteSummary returns a one-line summary of delete results.
func FormatDeleteSummary(s DeleteSummary) string {
	return fmt.Sprintf("%d deleted, %d skipped", s.Deleted, s.Skipped)
}
