package envfile

import (
	"fmt"
	"strings"
)

// FormatTransformResults formats the transform results for display.
func FormatTransformResults(results []TransformResult, colorize bool) string {
	var sb strings.Builder

	for _, r := range results {
		if r.Skipped {
			continue
		}
		if !r.Changed {
			continue
		}

		oldDisplay := r.OldValue
		newDisplay := r.NewValue
		if IsSensitive(r.Key) {
			oldDisplay = "[redacted]"
			newDisplay = "[redacted]"
		}

		line := fmt.Sprintf("  ~ %s: %q -> %q", r.Key, oldDisplay, newDisplay)
		if colorize {
			line = "\033[33m" + line + "\033[0m"
		}
		sb.WriteString(line + "\n")
	}

	return sb.String()
}

// FormatTransformSummary returns a one-line summary of transform results.
func FormatTransformSummary(results []TransformResult) string {
	changed, skipped := 0, 0
	for _, r := range results {
		if r.Skipped {
			skipped++
		} else if r.Changed {
			changed++
		}
	}
	parts := []string{fmt.Sprintf("%d transformed", changed)}
	if skipped > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", skipped))
	}
	return strings.Join(parts, ", ")
}

// FormatSkippedTransforms returns lines for skipped entries with reasons.
func FormatSkippedTransforms(results []TransformResult) string {
	var sb strings.Builder
	for _, r := range results {
		if !r.Skipped {
			continue
		}
		sb.WriteString(fmt.Sprintf("  - %s (%s)\n", r.Key, r.Reason))
	}
	return sb.String()
}
