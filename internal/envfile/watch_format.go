package envfile

import (
	"fmt"
	"strings"
)

// FormatWatchEvent returns a human-readable description of a WatchEvent.
func FormatWatchEvent(e WatchEvent, colorize bool) string {
	timestamp := e.ChangedAt.Format("15:04:05")

	var sb strings.Builder
	if colorize {
		sb.WriteString("\033[33m")
	}
	sb.WriteString(fmt.Sprintf("[%s] %s changed", timestamp, e.Path))
	if colorize {
		sb.WriteString("\033[0m")
	}
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("  old: %s\n", shortSum(e.OldSum)))
	sb.WriteString(fmt.Sprintf("  new: %s\n", shortSum(e.NewSum)))
	return sb.String()
}

// FormatWatchSummary returns a one-line summary of how many changes were seen.
func FormatWatchSummary(path string, count int) string {
	if count == 0 {
		return fmt.Sprintf("watch: no changes detected in %s", path)
	}
	plural := "change"
	if count != 1 {
		plural = "changes"
	}
	return fmt.Sprintf("watch: %d %s detected in %s", count, plural, path)
}

// shortSum returns the first 12 characters of a hex checksum.
func shortSum(sum string) string {
	if len(sum) <= 12 {
		return sum
	}
	return sum[:12] + "..."
}
