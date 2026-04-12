package envfile

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ExportFormat represents the supported export formats.
type ExportFormat string

const (
	FormatDotEnv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
	FormatShell  ExportFormat = "shell"
)

// Export serializes a slice of Entry values into the given format.
func Export(entries []Entry, format ExportFormat, redact bool) (string, error) {
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	switch format {
	case FormatDotEnv:
		return exportDotEnv(sorted, redact), nil
	case FormatJSON:
		return exportJSON(sorted, redact)
	case FormatShell:
		return exportShell(sorted, redact), nil
	default:
		return "", fmt.Errorf("unsupported export format: %q", format)
	}
}

func exportDotEnv(entries []Entry, redact bool) string {
	var sb strings.Builder
	for _, e := range entries {
		val := e.Value
		if redact && IsSensitive(e.Key) {
			val = "***"
		}
		if needsQuoting(val) {
			val = `"` + val + `"`
		}
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, val)
	}
	return sb.String()
}

func exportJSON(entries []Entry, redact bool) (string, error) {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		val := e.Value
		if redact && IsSensitive(e.Key) {
			val = "***"
		}
		m[e.Key] = val
	}
	out, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(out) + "\n", nil
}

func exportShell(entries []Entry, redact bool) string {
	var sb strings.Builder
	for _, e := range entries {
		val := e.Value
		if redact && IsSensitive(e.Key) {
			val = "***"
		}
		fmt.Fprintf(&sb, "export %s=%q\n", e.Key, val)
	}
	return sb.String()
}
