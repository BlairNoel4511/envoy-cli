package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatLockSummary returns a human-readable summary of a LockFile.
func FormatLockSummary(lf *LockFile) string {
	if lf == nil || len(lf.Entries) == 0 {
		return "No pinned entries."
	}

	keys := make([]string, 0, len(lf.Entries))
	for k := range lf.Entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Pinned entries (%d):\n", len(keys)))
	for _, k := range keys {
		e := lf.Entries[k]
		val := e.Value
		if IsSensitive(k) {
			val = "[REDACTED]"
		}
		comment := ""
		if e.Comment != "" {
			comment = fmt.Sprintf(" # %s", e.Comment)
		}
		sb.WriteString(fmt.Sprintf("  %s=%s  (by %s on %s)%s\n",
			k, val, e.PinnedBy, e.PinnedAt.Format("2006-01-02"), comment))
	}
	return sb.String()
}

// FormatApplyResult returns a human-readable summary of a LockApplyResult.
func FormatApplyResult(r LockApplyResult) string {
	var sb strings.Builder
	if len(r.Pinned) > 0 {
		sort.Strings(r.Pinned)
		sb.WriteString(fmt.Sprintf("Pinned: %s\n", strings.Join(r.Pinned, ", ")))
	}
	if len(r.Skipped) > 0 {
		sort.Strings(r.Skipped)
		sb.WriteString(fmt.Sprintf("Skipped (identical): %s\n", strings.Join(r.Skipped, ", ")))
	}
	if len(r.Conflict) > 0 {
		sort.Strings(r.Conflict)
		sb.WriteString(fmt.Sprintf("Conflicts (not overwritten): %s\n", strings.Join(r.Conflict, ", ")))
	}
	if sb.Len() == 0 {
		return "No changes applied."
	}
	return sb.String()
}
