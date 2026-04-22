package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatDiffApplyResults returns a human-readable summary of applied diff results.
func FormatDiffApplyResults(results []DiffApplyResult, colorize bool) string {
	if len(results) == 0 {
		return "(no changes)"
	}

	sorted := make([]DiffApplyResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	var sb strings.Builder
	for _, r := range sorted {
		line := formatDiffApplyLine(r, colorize)
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

func formatDiffApplyLine(r DiffApplyResult, colorize bool) string {
	val := r.NewValue
	if r.Sensitive {
		val = "[redacted]"
	}

	switch r.Action {
	case "added":
		msg := fmt.Sprintf("+ %s = %s", r.Key, val)
		if colorize {
			return "\033[32m" + msg + "\033[0m"
		}
		return msg
	case "updated":
		oldVal := r.OldValue
		if r.Sensitive {
			oldVal = "[redacted]"
		}
		msg := fmt.Sprintf("~ %s: %s -> %s", r.Key, oldVal, val)
		if colorize {
			return "\033[33m" + msg + "\033[0m"
		}
		return msg
	case "removed":
		msg := fmt.Sprintf("- %s", r.Key)
		if colorize {
			return "\033[31m" + msg + "\033[0m"
		}
		return msg
	case "skipped":
		msg := fmt.Sprintf("  %s (skipped)", r.Key)
		if colorize {
			return "\033[90m" + msg + "\033[0m"
		}
		return msg
	}
	return ""
}

// FormatDiffApplySummary returns a one-line summary of the apply operation.
func FormatDiffApplySummary(s DiffApplySummary) string {
	return fmt.Sprintf("applied: +%d added, ~%d updated, -%d removed, %d skipped",
		s.Added, s.Updated, s.Removed, s.Skipped)
}
