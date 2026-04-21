package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatCommentResults returns a human-readable string for a slice of CommentResults.
func FormatCommentResults(results []CommentResult, colorize bool) string {
	if len(results) == 0 {
		return ""
	}
	sorted := make([]CommentResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})
	var sb strings.Builder
	for _, r := range sorted {
		sb.WriteString(formatCommentLine(r, colorize))
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatCommentLine(r CommentResult, colorize bool) string {
	switch r.Status {
	case "added":
		line := fmt.Sprintf("+ %s # %s", r.Key, r.Comment)
		if colorize {
			return "\033[32m" + line + "\033[0m"
		}
		return line
	case "updated":
		line := fmt.Sprintf("~ %s # %s", r.Key, r.Comment)
		if colorize {
			return "\033[33m" + line + "\033[0m"
		}
		return line
	case "not_found":
		line := fmt.Sprintf("! %s (not found)", r.Key)
		if colorize {
			return "\033[31m" + line + "\033[0m"
		}
		return line
	default:
		return fmt.Sprintf("  %s (unchanged)", r.Key)
	}
}

// FormatCommentSummary returns a one-line summary of a CommentSummary.
func FormatCommentSummary(sum CommentSummary) string {
	return fmt.Sprintf("comment: %d added, %d updated, %d unchanged, %d not found",
		sum.Added, sum.Updated, sum.Unchanged, sum.NotFound)
}
