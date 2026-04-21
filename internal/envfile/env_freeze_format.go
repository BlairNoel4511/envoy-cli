package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatFreezeResults returns a human-readable summary of freeze operation results.
// Each line indicates whether a key was frozen, already frozen, skipped, or not found.
func FormatFreezeResults(results []FreezeResult, colorize bool) string {
	if len(results) == 0 {
		return "no keys processed"
	}

	sorted := make([]FreezeResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	var sb strings.Builder
	for _, r := range sorted {
		sb.WriteString(formatFreezeLine(r, colorize))
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// formatFreezeLine formats a single FreezeResult into a display line.
func formatFreezeLine(r FreezeResult, colorize bool) string {
	switch r.Status {
	case FreezeStatusFrozen:
		line := fmt.Sprintf("[frozen]  %s", r.Key)
		if colorize {
			return "\033[32m" + line + "\033[0m"
		}
		return line
	case FreezeStatusAlreadyFrozen:
		line := fmt.Sprintf("[skipped] %s (already frozen)", r.Key)
		if colorize {
			return "\033[33m" + line + "\033[0m"
		}
		return line
	case FreezeStatusRefrozen:
		line := fmt.Sprintf("[updated] %s (re-frozen)", r.Key)
		if colorize {
			return "\033[36m" + line + "\033[0m"
		}
		return line
	case FreezeStatusNotFound:
		line := fmt.Sprintf("[missing] %s", r.Key)
		if colorize {
			return "\033[31m" + line + "\033[0m"
		}
		return line
	default:
		return fmt.Sprintf("[unknown] %s", r.Key)
	}
}

// FormatFrozenList returns a formatted list of all currently frozen keys with
// their associated environments and optional reasons.
func FormatFrozenList(entries []FrozenEntry, colorize bool) string {
	if len(entries) == 0 {
		return "no frozen keys"
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	var sb strings.Builder
	for _, e := range entries {
		line := fmt.Sprintf("  %s", e.Key)
		if e.Env != "" {
			line += fmt.Sprintf(" [env:%s]", e.Env)
		}
		if e.Reason != "" {
			line += fmt.Sprintf(" — %s", e.Reason)
		}
		if colorize {
			line = "\033[36m" + line + "\033[0m"
		}
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// FormatFreezeSummary returns a one-line summary of freeze operation counts.
func FormatFreezeSummary(results []FreezeResult) string {
	var frozen, skipped, refrozen, missing int
	for _, r := range results {
		switch r.Status {
		case FreezeStatusFrozen:
			frozen++
		case FreezeStatusAlreadyFrozen:
			skipped++
		case FreezeStatusRefrozen:
			refrozen++
		case FreezeStatusNotFound:
			missing++
		}
	}
	parts := []string{}
	if frozen > 0 {
		parts = append(parts, fmt.Sprintf("%d frozen", frozen))
	}
	if refrozen > 0 {
		parts = append(parts, fmt.Sprintf("%d re-frozen", refrozen))
	}
	if skipped > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", skipped))
	}
	if missing > 0 {
		parts = append(parts, fmt.Sprintf("%d not found", missing))
	}
	if len(parts) == 0 {
		return "nothing to freeze"
	}
	return strings.Join(parts, ", ")
}
