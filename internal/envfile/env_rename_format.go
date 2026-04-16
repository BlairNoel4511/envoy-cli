package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatRenameResults returns a human-readable summary of rename operations.
func FormatRenameResults(results []RenameResult, colorize bool) string {
	if len(results) == 0 {
		return "no rename operations performed"
	}

	sorted := make([]RenameResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].OldKey < sorted[j].OldKey
	})

	var sb strings.Builder
	for _, r := range sorted {
		sb.WriteString(formatRenameLine(r, colorize))
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

func formatRenameLine(r RenameResult, colorize bool) string {
	switch r.Status {
	case "renamed":
		line := fmt.Sprintf("  renamed  %s → %s", r.OldKey, r.NewKey)
		if colorize {
			return "\033[32m" + line + "\033[0m"
		}
		return line
	case "conflict":
		line := fmt.Sprintf("  conflict %s → %s (%s)", r.OldKey, r.NewKey, r.Comment)
		if colorize {
			return "\033[33m" + line + "\033[0m"
		}
		return line
	case "not_found":
		line := fmt.Sprintf("  missing  %s (key not found)", r.OldKey)
		if colorize {
			return "\033[31m" + line + "\033[0m"
		}
		return line
	default:
		line := fmt.Sprintf("  skipped  %s → %s", r.OldKey, r.NewKey)
		if colorize {
			return "\033[90m" + line + "\033[0m"
		}
		return line
	}
}

// FormatRenameSummaryLine returns a one-line summary string.
func FormatRenameSummaryLine(results []RenameResult) string {
	sm := RenameSummary(results)
	return fmt.Sprintf("renamed: %d  conflict: %d  not_found: %d  skipped: %d",
		sm["renamed"], sm["conflict"], sm["not_found"], sm["skipped"])
}
