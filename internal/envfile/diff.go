package envfile

// DiffType represents the type of change between two env files.
type DiffType string

const (
	DiffAdded   DiffType = "added"
	DiffRemoved DiffType = "removed"
	DiffChanged DiffType = "changed"
	DiffUnchanged DiffType = "unchanged"
)

// DiffEntry represents a single key difference between two env maps.
type DiffEntry struct {
	Key      string
	OldValue string
	NewValue string
	Type     DiffType
}

// Diff compares two env maps and returns a slice of DiffEntry describing
// what was added, removed, changed, or unchanged between local and remote.
func Diff(local, remote map[string]string) []DiffEntry {
	var entries []DiffEntry
	seen := make(map[string]bool)

	for k, localVal := range local {
		seen[k] = true
		remoteVal, exists := remote[k]
		switch {
		case !exists:
			entries = append(entries, DiffEntry{Key: k, OldValue: localVal, NewValue: "", Type: DiffRemoved})
		case localVal != remoteVal:
			entries = append(entries, DiffEntry{Key: k, OldValue: localVal, NewValue: remoteVal, Type: DiffChanged})
		default:
			entries = append(entries, DiffEntry{Key: k, OldValue: localVal, NewValue: remoteVal, Type: DiffUnchanged})
		}
	}

	for k, remoteVal := range remote {
		if !seen[k] {
			entries = append(entries, DiffEntry{Key: k, OldValue: "", NewValue: remoteVal, Type: DiffAdded})
		}
	}

	return entries
}

// HasChanges returns true if the diff contains any added, removed, or changed entries.
func HasChanges(entries []DiffEntry) bool {
	for _, e := range entries {
		if e.Type != DiffUnchanged {
			return true
		}
	}
	return false
}
