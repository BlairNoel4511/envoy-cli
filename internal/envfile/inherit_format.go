package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatInheritResult returns a detailed formatted string of the inherit result.
func FormatInheritResult(r InheritResult, colorize bool) string {
	var sb strings.Builder

	green := func(s string) string {
		if colorize {
			return "\033[32m" + s + "\033[0m"
		}
		return s
	}
	yellow := func(s string) string {
		if colorize {
			return "\033[33m" + s + "\033[0m"
		}
		return s
	}
	gray := func(s string) string {
		if colorize {
			return "\033[90m" + s + "\033[0m"
		}
		return s
	}

	added := make([]string, len(r.Added))
	copy(added, r.Added)
	sort.Strings(added)
	for _, k := range added {
		sb.WriteString(green(fmt.Sprintf("  + %s (inherited)\n", k)))
	}

	overwritten := make([]string, len(r.Overwritten))
	copy(overwritten, r.Overwritten)
	sort.Strings(overwritten)
	for _, k := range overwritten {
		sb.WriteString(yellow(fmt.Sprintf("  ~ %s (overwritten)\n", k)))
	}

	skipped := make([]string, len(r.Skipped))
	copy(skipped, r.Skipped)
	sort.Strings(skipped)
	for _, k := range skipped {
		sb.WriteString(gray(fmt.Sprintf("  - %s (skipped)\n", k)))
	}

	sb.WriteString(fmt.Sprintf("\nSummary: %s\n", InheritSummary(r)))
	return sb.String()
}
