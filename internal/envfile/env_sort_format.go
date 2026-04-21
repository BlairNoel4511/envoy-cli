package envfile

import (
	"fmt"
	"strings"
)

// FormatSortResults returns a formatted string showing before/after order.
func FormatSortResults(r SortResult, colorize bool) string {
	if len(r.Entries) == 0 {
		return "(no entries)"
	}

	origIndex := make(map[string]int, len(r.Original))
	for i, e := range r.Original {
		origIndex[e.Key] = i
	}

	var sb strings.Builder
	for newPos, e := range r.Entries {
		oldPos := origIndex[e.Key]
		moved := oldPos != newPos

		line := fmt.Sprintf("  %-30s", e.Key)
		if moved {
			annotation := fmt.Sprintf(" (was #%d, now #%d)", oldPos+1, newPos+1)
			if colorize {
				line = "\033[33m" + line + annotation + "\033[0m"
			} else {
				line = line + annotation
			}
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}

// FormatSortSummaryLine returns a concise one-line summary.
func FormatSortSummaryLine(r SortResult, colorize bool) string {
	summary := SortSummary(r)
	if colorize && r.Reordered > 0 {
		return "\033[32m" + summary + "\033[0m"
	}
	return summary
}
