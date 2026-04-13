package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatDeprecationList renders all deprecated entries as a human-readable list.
func FormatDeprecationList(store *DeprecationStore, colorize bool) string {
	if len(store.Entries) == 0 {
		return "no deprecated keys\n"
	}

	entries := store.List()
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	var sb strings.Builder
	for _, e := range entries {
		line := formatDeprecatedEntry(e, colorize)
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}

// FormatDeprecationWarnings formats only the entries found in a set of live entries.
func FormatDeprecationWarnings(hits []DeprecatedEntry, colorize bool) string {
	if len(hits) == 0 {
		return ""
	}
	sort.Slice(hits, func(i, j int) bool {
		return hits[i].Key < hits[j].Key
	})
	var sb strings.Builder
	sb.WriteString("deprecated keys in use:\n")
	for _, e := range hits {
		sb.WriteString("  ")
		sb.WriteString(formatDeprecatedEntry(e, colorize))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func formatDeprecatedEntry(e DeprecatedEntry, colorize bool) string {
	status := string(e.Status)
	if colorize {
		switch e.Status {
		case DeprecationRemoved:
			status = "\033[31m" + status + "\033[0m"
		case DeprecationWarning:
			status = "\033[33m" + status + "\033[0m"
		case DeprecationActive:
			status = "\033[32m" + status + "\033[0m"
		}
	}
	line := fmt.Sprintf("[%s] %s — %s", status, e.Key, e.Reason)
	if e.ReplacedBy != "" {
		line += fmt.Sprintf(" (use %s instead)", e.ReplacedBy)
	}
	return line
}

// FormatDeprecationSummary returns a one-line summary count.
func FormatDeprecationSummary(store *DeprecationStore) string {
	return fmt.Sprintf("%d deprecated key(s) tracked", len(store.Entries))
}
