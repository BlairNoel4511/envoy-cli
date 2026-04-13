package envfile

import (
	"fmt"
	"strings"
)

// FormatResolveResults returns a human-readable summary of resolution results.
func FormatResolveResults(results []ResolveResult, colorize bool) string {
	if len(results) == 0 {
		return "no entries resolved\n"
	}
	var b strings.Builder
	for _, r := range results {
		if !r.Changed && len(r.Unresolved) == 0 {
			continue
		}
		if r.Changed {
			line := fmt.Sprintf("~ %s: %q -> %q", r.Key, r.Original, r.Resolved)
			if colorize {
				line = "\033[33m" + line + "\033[0m"
			}
			b.WriteString(line + "\n")
		}
		for _, u := range r.Unresolved {
			line := fmt.Sprintf("  ! unresolved variable: $%s", u)
			if colorize {
				line = "\033[31m" + line + "\033[0m"
			}
			b.WriteString(line + "\n")
		}
	}
	if b.Len() == 0 {
		return "all variables resolved without changes\n"
	}
	return b.String()
}

// FormatResolveSummary returns a one-line summary.
func FormatResolveSummary(results []ResolveResult) string {
	changed, unresolved := 0, 0
	for _, r := range results {
		if r.Changed {
			changed++
		}
		unresolved += len(r.Unresolved)
	}
	return fmt.Sprintf("resolved: %d changed, %d unresolved variable(s)", changed, unresolved)
}
