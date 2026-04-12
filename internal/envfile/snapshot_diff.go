package envfile

import "fmt"

// SnapshotDiff compares two snapshots and returns a DiffResult.
// It reuses the existing Diff function by converting snapshots to maps.
func SnapshotDiff(before, after Snapshot) DiffResult {
	return Diff(before.ToMap(), after.ToMap())
}

// DiffResult holds the categorised changes between two env states.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // [before, after]
}

// SnapshotSummary returns a human-readable summary of changes between snapshots.
func SnapshotSummary(before, after Snapshot) string {
	result := SnapshotDiff(before, after)

	if !HasChanges(result) {
		return fmt.Sprintf("No changes between snapshots (%s → %s)",
			before.Timestamp.Format("2006-01-02 15:04:05"),
			after.Timestamp.Format("2006-01-02 15:04:05"),
		)
	}

	summary := fmt.Sprintf("Snapshot diff (%s → %s):\n",
		before.Timestamp.Format("2006-01-02 15:04:05"),
		after.Timestamp.Format("2006-01-02 15:04:05"),
	)

	for k, v := range result.Added {
		display := maskIfSensitive(k, v)
		summary += fmt.Sprintf("  + %s=%s\n", k, display)
	}
	for k := range result.Removed {
		summary += fmt.Sprintf("  - %s\n", k)
	}
	for k, pair := range result.Changed {
		before := maskIfSensitive(k, pair[0])
		after := maskIfSensitive(k, pair[1])
		summary += fmt.Sprintf("  ~ %s: %s → %s\n", k, before, after)
	}

	return summary
}

// ChecksumChanged reports whether two snapshots have different checksums.
func ChecksumChanged(before, after Snapshot) bool {
	return before.Checksum != after.Checksum
}
