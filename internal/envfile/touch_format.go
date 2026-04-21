package envfile

import (
	"fmt"
	"sort"
	"strings"
)

const touchTimeLayout = "2006-01-02 15:04:05 UTC"

// FormatTouchResults returns a human-readable table of touch results.
func FormatTouchResults(results []TouchResult, colorize bool) string {
	if len(results) == 0 {
		return "(no keys touched)\n"
	}

	sorted := make([]TouchResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	var sb strings.Builder
	for _, r := range sorted {
		switch {
		case r.Skipped:
			line := fmt.Sprintf("  ~ %-30s skipped (%s)\n", r.Key, r.Reason)
			if colorize {
				line = "\033[33m" + line + "\033[0m"
			}
			sb.WriteString(line)
		case r.WasSet && !r.PrevTouch.IsZero():
			line := fmt.Sprintf("  ↻ %-30s updated → %s\n", r.Key, r.NewTouch.Format(touchTimeLayout))
			if colorize {
				line = "\033[36m" + line + "\033[0m"
			}
			sb.WriteString(line)
		default:
			line := fmt.Sprintf("  + %-30s set     → %s\n", r.Key, r.NewTouch.Format(touchTimeLayout))
			if colorize {
				line = "\033[32m" + line + "\033[0m"
			}
			sb.WriteString(line)
		}
	}
	return sb.String()
}

// FormatTouchSummary returns a one-line summary of touch results.
func FormatTouchSummary(s TouchSummary) string {
	return fmt.Sprintf("touched %d key(s), skipped %d", s.Touched, s.Skipped)
}
