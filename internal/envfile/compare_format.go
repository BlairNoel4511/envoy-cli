package envfile

import (
	"fmt"
	"strings"
)

// FormatCompareResult returns a human-readable diff-style string for a CompareResult.
// Sensitive values are redacted when redact is true.
func FormatCompareResult(r CompareResult, redact bool) string {
	var sb strings.Builder

	for _, e := range r.OnlyInRight {
		val := maskCompareValue(e.Key, e.Value, redact)
		sb.WriteString(fmt.Sprintf("+ %s=%s\n", e.Key, val))
	}

	for _, e := range r.OnlyInLeft {
		val := maskCompareValue(e.Key, e.Value, redact)
		sb.WriteString(fmt.Sprintf("- %s=%s\n", e.Key, val))
	}

	for _, p := range r.Different {
		lv := maskCompareValue(p.Key, p.LeftValue, redact)
		rv := maskCompareValue(p.Key, p.RightValue, redact)
		sb.WriteString(fmt.Sprintf("~ %s: %s -> %s\n", p.Key, lv, rv))
	}

	for _, e := range r.Identical {
		val := maskCompareValue(e.Key, e.Value, redact)
		sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, val))
	}

	return sb.String()
}

// FormatCompareSummaryLine returns a styled summary line for CLI output.
func FormatCompareSummaryLine(r CompareResult) string {
	parts := []string{}
	if len(r.OnlyInRight) > 0 {
		parts = append(parts, fmt.Sprintf("%d added", len(r.OnlyInRight)))
	}
	if len(r.OnlyInLeft) > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", len(r.OnlyInLeft)))
	}
	if len(r.Different) > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", len(r.Different)))
	}
	if len(parts) == 0 {
		return "No differences found."
	}
	return strings.Join(parts, ", ") + "."
}

func maskCompareValue(key, value string, redact bool) string {
	if redact && IsSensitive(key) {
		return "***"
	}
	return value
}
