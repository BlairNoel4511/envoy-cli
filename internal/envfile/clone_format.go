package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatCloneDiff renders a colourised, sorted diff-style view of a clone
// operation showing which keys were cloned and which were skipped.
func FormatCloneDiff(result CloneResult, colorize bool) string {
	var sb strings.Builder

	cloned := make([]string, len(result.Cloned))
	copy(cloned, result.Cloned)
	sort.Strings(cloned)

	skipped := make([]string, len(result.Skipped))
	copy(skipped, result.Skipped)
	sort.Strings(skipped)

	for _, k := range cloned {
		line := fmt.Sprintf("+ %s", k)
		if colorize {
			line = "\033[32m" + line + "\033[0m"
		}
		sb.WriteString(line + "\n")
	}

	for _, k := range skipped {
		line := fmt.Sprintf("~ %s (skipped)", k)
		if colorize {
			line = "\033[33m" + line + "\033[0m"
		}
		sb.WriteString(line + "\n")
	}

	for _, e := range result.Errors {
		line := fmt.Sprintf("! %s", e)
		if colorize {
			line = "\033[31m" + line + "\033[0m"
		}
		sb.WriteString(line + "\n")
	}

	sb.WriteString(fmt.Sprintf("\nSummary: %d cloned, %d skipped, %d errors\n",
		len(result.Cloned), len(result.Skipped), len(result.Errors)))

	return sb.String()
}

// CloneResultHasChanges returns true if any keys were cloned.
func CloneResultHasChanges(r CloneResult) bool {
	return len(r.Cloned) > 0
}
