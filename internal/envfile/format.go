package envfile

import (
	"fmt"
	"sort"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

// FormatDiff returns a human-readable, colored string representation of diff entries.
// Sensitive values are redacted automatically.
func FormatDiff(entries []DiffEntry, colorize bool) string {
	sorted := make([]DiffEntry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	var sb strings.Builder
	for _, e := range sorted {
		oldVal := maskIfSensitive(e.Key, e.OldValue)
		newVal := maskIfSensitive(e.Key, e.NewValue)

		var line string
		switch e.Type {
		case DiffAdded:
			line = fmt.Sprintf("+ %s=%s", e.Key, newVal)
			if colorize {
				line = colorGreen + line + colorReset
			}
		case DiffRemoved:
			line = fmt.Sprintf("- %s=%s", e.Key, oldVal)
			if colorize {
				line = colorRed + line + colorReset
			}
		case DiffChanged:
			line = fmt.Sprintf("~ %s: %s → %s", e.Key, oldVal, newVal)
			if colorize {
				line = colorYellow + line + colorReset
			}
		case DiffUnchanged:
			line = fmt.Sprintf("  %s=%s", e.Key, oldVal)
			if colorize {
				line = colorGray + line + colorReset
			}
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}

func maskIfSensitive(key, value string) string {
	if IsSensitive(key) && value != "" {
		return "****"
	}
	return value
}
