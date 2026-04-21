package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatProtectResults returns a human-readable list of protect operation results.
func FormatProtectResults(results []ProtectResult, colorize bool) string {
	if len(results) == 0 {
		return "no keys processed"
	}

	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(formatProtectLine(r, colorize))
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

func formatProtectLine(r ProtectResult, colorize bool) string {
	switch {
	case r.Protected:
		if colorize {
			return fmt.Sprintf("\033[32m+ protected\033[0m  %s", r.Key)
		}
		return fmt.Sprintf("+ protected  %s", r.Key)
	case r.Already:
		if colorize {
			return fmt.Sprintf("\033[33m~ already\033[0m    %s", r.Key)
		}
		return fmt.Sprintf("~ already    %s", r.Key)
	case r.NotFound:
		if colorize {
			return fmt.Sprintf("\033[31m! not found\033[0m  %s", r.Key)
		}
		return fmt.Sprintf("! not found  %s", r.Key)
	}
	return r.Key
}

// FormatProtectedList returns a sorted list of all currently protected keys.
func FormatProtectedList(store *ProtectStore) string {
	keys := store.All()
	if len(keys) == 0 {
		return "no protected keys"
	}
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("  🔒 %s\n", k))
	}
	return strings.TrimRight(sb.String(), "\n")
}
