package envfile

import (
	"fmt"
	"strings"
)

// FormatSetResults returns a formatted, optionally colorized list of set operation results.
func FormatSetResults(results []SetResult, colorize bool) string {
	if len(results) == 0 {
		return "(no changes)"
	}

	var sb strings.Builder
	for _, r := range results {
		line := formatSetLine(r, colorize)
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// SummarizeSetResults returns a short summary string of how many keys were
// added, updated, skipped, and unchanged from the given results.
func SummarizeSetResults(results []SetResult) string {
	counts := make(map[string]int)
	for _, r := range results {
		counts[r.Action]++
	}
	return fmt.Sprintf("added=%d updated=%d skipped=%d unchanged=%d",
		counts["added"], counts["updated"], counts["skipped"], counts["unchanged"])
}

func formatSetLine(r SetResult, colorize bool) string {
	masked := maskIfSensitive(r.Key, r.Value)
	switch r.Action {
	case "added":
		msg := fmt.Sprintf("+ %s=%s", r.Key, masked)
		if colorize {
			return "\033[32m" + msg + "\033[0m"
		}
		return msg
	case "updated":
		oldMasked := maskIfSensitive(r.Key, r.OldValue)
		msg := fmt.Sprintf("~ %s: %s → %s", r.Key, oldMasked, masked)
		if colorize {
			return "\033[33m" + msg + "\033[0m"
		}
		return msg
	case "skipped":
		msg := fmt.Sprintf("! %s (skipped, key exists)", r.Key)
		if colorize {
			return "\033[31m" + msg + "\033[0m"
		}
		return msg
	case "unchanged":
		msg := fmt.Sprintf("= %s (unchanged)", r.Key)
		if colorize {
			return "\033[90m" + msg + "\033[0m"
		}
		return msg
	}
	return fmt.Sprintf("? %s", r.Key)
}
