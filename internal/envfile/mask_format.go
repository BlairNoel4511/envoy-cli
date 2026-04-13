package envfile

import (
	"fmt"
	"strings"
)

// FormatMaskResults formats a slice of MaskResult for display.
// If colorize is true, ANSI codes are added for masked entries.
func FormatMaskResults(results []MaskResult, colorize bool) string {
	if len(results) == 0 {
		return "(no entries)"
	}
	var sb strings.Builder
	for _, r := range results {
		var line string
		if r.WasMasked {
			if colorize {
				line = fmt.Sprintf("\033[33m%s\033[0m=%s", r.Key, r.Masked)
			} else {
				line = fmt.Sprintf("%s=%s", r.Key, r.Masked)
			}
		} else {
			line = fmt.Sprintf("%s=%s", r.Key, r.Masked)
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}

// FormatMaskSummary returns a one-line summary of masking results.
func FormatMaskSummary(results []MaskResult) string {
	total := len(results)
	masked := 0
	for _, r := range results {
		if r.WasMasked {
			masked++
		}
	}
	return fmt.Sprintf("%d/%d values masked", masked, total)
}

// FormatMaskedEntry formats a single MaskResult as a key=value line.
func FormatMaskedEntry(r MaskResult, colorize bool) string {
	if r.WasMasked && colorize {
		return fmt.Sprintf("\033[33m%s\033[0m=%s", r.Key, r.Masked)
	}
	return fmt.Sprintf("%s=%s", r.Key, r.Masked)
}
