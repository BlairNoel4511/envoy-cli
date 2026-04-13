package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatRollbackPlan returns a human-readable summary of a RollbackPlan.
func FormatRollbackPlan(plan RollbackPlan, colorize bool) string {
	if len(plan.Entries) == 0 {
		return fmt.Sprintf("rollback plan %q: no changes to revert\n", plan.OperationID)
	}

	entries := make([]RollbackEntry, len(plan.Entries))
	copy(entries, plan.Entries)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("rollback plan %q (%d change(s)):\n", plan.OperationID, len(entries)))

	for _, e := range entries {
		oldDisplay := maskIfSensitive(e.Key, e.OldValue)
		newDisplay := maskIfSensitive(e.Key, e.NewValue)

		var line string
		if !e.HadKey {
			line = fmt.Sprintf("  - %s (remove, was added with value %s)", e.Key, newDisplay)
		} else {
			line = fmt.Sprintf("  ~ %s: %s -> %s", e.Key, newDisplay, oldDisplay)
		}

		if colorize {
			if !e.HadKey {
				line = "\033[31m" + line + "\033[0m"
			} else {
				line = "\033[33m" + line + "\033[0m"
			}
		}
		sb.WriteString(line + "\n")
	}

	return sb.String()
}

// FormatRollbackSummary returns a one-line summary after applying a rollback.
func FormatRollbackSummary(plan RollbackPlan) string {
	removed := 0
	reverted := 0
	for _, e := range plan.Entries {
		if e.HadKey {
			reverted++
		} else {
			removed++
		}
	}
	return fmt.Sprintf("rollback applied: %d key(s) reverted, %d key(s) removed", reverted, removed)
}
