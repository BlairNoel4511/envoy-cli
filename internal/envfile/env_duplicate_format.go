package envfile

import (
	"fmt"
	"strings"
)

// FormatDuplicateResult formats a single DuplicateResult for display.
func FormatDuplicateResult(r DuplicateResult, colorize bool) string {
	var sb strings.Builder

	val := r.Value
	if r.Sensitive {
		val = "[redacted]"
	}

	switch r.Status {
	case "duplicated":
		line := fmt.Sprintf("+ %s => %s = %s", r.SourceKey, r.DestKey, val)
		if colorize {
			line = "\033[32m" + line + "\033[0m"
		}
		sb.WriteString(line)
	case "skipped":
		line := fmt.Sprintf("~ %s => %s (skipped)", r.SourceKey, r.DestKey)
		if colorize {
			line = "\033[33m" + line + "\033[0m"
		}
		sb.WriteString(line)
	case "unchanged":
		line := fmt.Sprintf("= %s => %s (unchanged)", r.SourceKey, r.DestKey)
		sb.WriteString(line)
	case "source_not_found":
		line := fmt.Sprintf("! %s (source not found)", r.SourceKey)
		if colorize {
			line = "\033[31m" + line + "\033[0m"
		}
		sb.WriteString(line)
	}

	return sb.String()
}

// FormatDuplicateSummary returns a human-readable summary line.
func FormatDuplicateSummary(s DuplicateSummary) string {
	return fmt.Sprintf(
		"duplicate: %d duplicated, %d skipped, %d unchanged, %d not found",
		s.Duplicated, s.Skipped, s.Unchanged, s.NotFound,
	)
}
