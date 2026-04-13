package envfile

import (
	"fmt"
	"strings"
)

// FormatPatchResults returns a human-readable summary of patch results.
func FormatPatchResults(results []PatchResult, colorize bool) string {
	var sb strings.Builder
	for _, r := range results {
		line := formatPatchLine(r, colorize)
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatPatchLine(r PatchResult, colorize bool) string {
	var status, color, reset string
	reset = ""
	if colorize {
		reset = "\033[0m"
	}

	switch {
	case r.Applied:
		status = "applied"
		if colorize {
			color = "\033[32m" // green
		}
	case r.Skipped:
		status = "skipped"
		if colorize {
			color = "\033[33m" // yellow
		}
	default:
		status = "unknown"
	}

	ops := string(r.Instruction.Op)
	key := r.Instruction.Key

	base := fmt.Sprintf("%s[%s] %s %s%s", color, status, ops, key, reset)
	if r.Skipped && r.Reason != "" {
		base += fmt.Sprintf(" (%s)", r.Reason)
	}
	if r.Instruction.Op == PatchSet && r.Applied {
		base += fmt.Sprintf(" = %s", r.Instruction.Value)
	}
	if r.Instruction.Op == PatchRename && r.Applied {
		base += fmt.Sprintf(" -> %s", r.Instruction.NewKey)
	}
	return base
}

// FormatPatchSummary returns a one-line summary of patch outcomes.
func FormatPatchSummary(results []PatchResult) string {
	applied, skipped := 0, 0
	for _, r := range results {
		if r.Applied {
			applied++
		} else if r.Skipped {
			skipped++
		}
	}
	return fmt.Sprintf("%d applied, %d skipped", applied, skipped)
}
