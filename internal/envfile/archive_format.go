package envfile

import (
	"fmt"
	"strings"
)

// FormatArchiveList returns a human-readable list of archive records.
func FormatArchiveList(records []ArchiveEntry, colorize bool) string {
	if len(records) == 0 {
		return "No archived snapshots found."
	}
	var sb strings.Builder
	for _, r := range records {
		line := fmt.Sprintf("[%s] %s — %d keys — %s",
			r.ID,
			r.CreatedAt.Format("2006-01-02 15:04:05"),
			len(r.Entries),
			r.Label,
		)
		if colorize {
			line = "\033[36m" + line + "\033[0m"
		}
		sb.WriteString(line + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// FormatArchiveDetail returns a formatted view of a single archive entry.
func FormatArchiveDetail(record ArchiveEntry, redact bool, colorize bool) string {
	var sb strings.Builder
	header := fmt.Sprintf("Archive: %s (%s)\nLabel: %s\nKeys: %d",
		record.ID,
		record.CreatedAt.Format("2006-01-02 15:04:05"),
		record.Label,
		len(record.Entries),
	)
	if colorize {
		header = "\033[1m" + header + "\033[0m"
	}
	sb.WriteString(header + "\n")
	for _, e := range record.Entries {
		val := e.Value
		if redact && IsSensitive(e.Key) {
			val = "***"
		}
		line := fmt.Sprintf("  %s=%s", e.Key, val)
		sb.WriteString(line + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// FormatArchiveSummary returns a one-line summary string.
func FormatArchiveSummary(count int) string {
	if count == 0 {
		return "Archive is empty."
	}
	return fmt.Sprintf("%d snapshot(s) stored in archive.", count)
}
