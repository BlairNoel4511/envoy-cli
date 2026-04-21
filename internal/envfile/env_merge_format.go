package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatMergeResults returns a human-readable, sorted list of merge actions.
func FormatMergeResults(results []MergeResult, colorize bool) string {
	if len(results) == 0 {
		return "(no changes)"
	}

	sorted := make([]MergeResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	var sb strings.Builder
	for _, r := range sorted {
		line := formatMergeLine(r, colorize)
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}

func formatMergeLine(r MergeResult, colorize bool) string {
	var prefix, color, reset string
	reset = ""
	switch r.Action {
	case "added":
		prefix = "+ "
		color = "\033[32m"
	case "updated":
		prefix = "~ "
		color = "\033[33m"
	case "skipped":
		prefix = "! "
		color = "\033[90m"
	case "unchanged":
		prefix = "  "
		color = ""
	}
	if colorize && color != "" {
		reset = "\033[0m"
	} else {
		color = ""
	}

	val := r.Value
	if r.Sensitive {
		val = "[redacted]"
	}
	return fmt.Sprintf("%s%s%s=%s%s", color, prefix, r.Key, val, reset)
}

// FormatMergeSummaryLine returns a one-line summary of a merge operation.
func FormatMergeSummaryLine(s MergeSummary) string {
	return fmt.Sprintf("merge: %d added, %d updated, %d skipped, %d unchanged",
		s.Added, s.Updated, s.Skipped, s.Unchanged)
}
