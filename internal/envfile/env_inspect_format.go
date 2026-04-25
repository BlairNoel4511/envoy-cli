package envfile

import (
	"fmt"
	"strings"
)

const (
	inspectColorReset  = "\033[0m"
	inspectColorYellow = "\033[33m"
	inspectColorRed    = "\033[31m"
	inspectColorGray   = "\033[90m"
)

// FormatInspectResults formats a slice of InspectResult for terminal output.
func FormatInspectResults(results []InspectResult, colorize bool) string {
	if len(results) == 0 {
		return "(no entries to inspect)\n"
	}
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(formatInspectLine(r, colorize))
	}
	return sb.String()
}

func formatInspectLine(r InspectResult, colorize bool) string {
	if !r.Found {
		if colorize {
			return fmt.Sprintf("%s%-30s%s  not found\n", inspectColorRed, r.Key, inspectColorReset)
		}
		return fmt.Sprintf("%-30s  not found\n", r.Key)
	}

	var flags []string
	if r.Sensitive {
		flags = append(flags, "sensitive")
	}
	if r.Redacted {
		flags = append(flags, "redacted")
	}
	if r.HasComment {
		flags = append(flags, fmt.Sprintf("comment=%q", r.Comment))
	}
	flagStr := ""
	if len(flags) > 0 {
		flagStr = "  [" + strings.Join(flags, ", ") + "]"
	}

	if colorize && r.Sensitive {
		return fmt.Sprintf("%s%-30s%s = %s%-40s%s  len=%d%s\n",
			inspectColorYellow, r.Key, inspectColorReset,
			inspectColorGray, r.Value, inspectColorReset,
			r.Length, flagStr)
	}
	return fmt.Sprintf("%-30s = %-40s  len=%d%s\n", r.Key, r.Value, r.Length, flagStr)
}

// FormatInspectSummary returns a formatted summary line.
func FormatInspectSummary(results []InspectResult) string {
	return InspectSummary(results) + "\n"
}
