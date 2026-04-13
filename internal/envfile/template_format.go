package envfile

import (
	"fmt"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
)

// FormatExpandedEntries returns a coloured, human-readable table of
// entries after template expansion, highlighting changed values.
func FormatExpandedEntries(original, expanded []Entry, colorize bool) string {
	orig := ToMap(original)
	exp := ToMap(expanded)

	var sb strings.Builder
	for _, e := range expanded {
		origVal := orig[e.Key]
		newVal := exp[e.Key]

		switch {
		case origVal == newVal:
			line := fmt.Sprintf("  %s=%s", e.Key, newVal)
			sb.WriteString(line + "\n")
		default:
			if colorize {
				sb.WriteString(fmt.Sprintf("%s~ %s=%s%s  (was: %s)\n",
					colorYellow, e.Key, newVal, colorReset, origVal))
			} else {
				sb.WriteString(fmt.Sprintf("~ %s=%s  (was: %s)\n", e.Key, newVal, origVal))
			}
		}
	}
	return sb.String()
}

// FormatMissingVars returns a warning block for unresolved template vars.
func FormatMissingVars(missing []string, colorize bool) string {
	if len(missing) == 0 {
		return ""
	}
	var sb strings.Builder
	if colorize {
		sb.WriteString(colorRed)
	}
	sb.WriteString(fmt.Sprintf("Warning: %d unresolved variable(s):\n", len(missing)))
	for _, k := range missing {
		sb.WriteString(fmt.Sprintf("  - %s\n", k))
	}
	if colorize {
		sb.WriteString(colorReset)
	}
	return sb.String()
}

// FormatExpandSummary produces a one-line summary of an expansion run.
func FormatExpandSummary(resolved, missing []string) string {
	return fmt.Sprintf("Expansion complete: %d resolved, %d missing.",
		len(resolved), len(missing))
}
