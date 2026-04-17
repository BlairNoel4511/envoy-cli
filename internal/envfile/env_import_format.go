package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatImportResults returns a human-readable summary of import results.
func FormatImportResults(results []ImportResult, colorize bool) string {
	if len(results) == 0 {
		return "  (nothing to import)\n"
	}

	sorted := make([]ImportResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	var sb strings.Builder
	for _, r := range sorted {
		sb.WriteString(formatImportLine(r, colorize))
	}
	return sb.String()
}

func formatImportLine(r ImportResult, colorize bool) string {
	val := r.Value
	if r.Sensitive {
		val = "***"
	}
	var prefix, reset string
	if colorize {
		reset = "\033[0m"
		switch r.Status {
		case "added":
			prefix = "\033[32m"
		case "overwritten":
			prefix = "\033[33m"
		case "skipped":
			prefix = "\033[90m"
		case "dry-run":
			prefix = "\033[36m"
		}
	}
	return fmt.Sprintf("  %s[%s]%s %s=%s\n", prefix, r.Status, reset, r.Key, val)
}

// FormatImportSummary returns a one-line summary of the import operation.
func FormatImportSummary(s ImportSummary) string {
	return fmt.Sprintf("import: %d added, %d overwritten, %d skipped (%d total)\n",
		s.Added, s.Overwritten, s.Skipped, s.Total)
}
