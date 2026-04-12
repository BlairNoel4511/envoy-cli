package envfile

import "fmt"

// LockApplyResult describes the outcome of applying a lock file to entries.
type LockApplyResult struct {
	Pinned   []string
	Skipped  []string
	Conflict []string
}

// ApplyLock enforces pinned values from a LockFile onto a slice of Entry.
// If overwrite is true, pinned values replace existing ones; otherwise conflicts are recorded.
func ApplyLock(entries []Entry, lf *LockFile, overwrite bool) ([]Entry, LockApplyResult, error) {
	if lf == nil {
		return nil, LockApplyResult{}, fmt.Errorf("lock file must not be nil")
	}

	result := LockApplyResult{}
	index := make(map[string]int, len(entries))
	for i, e := range entries {
		index[e.Key] = i
	}

	for key, locked := range lf.Entries {
		if i, exists := index[key]; exists {
			if entries[i].Value == locked.Value {
				result.Skipped = append(result.Skipped, key)
				continue
			}
			if overwrite {
				entries[i].Value = locked.Value
				result.Pinned = append(result.Pinned, key)
			} else {
				result.Conflict = append(result.Conflict, key)
			}
		} else {
			entries = append(entries, Entry{Key: key, Value: locked.Value})
			result.Pinned = append(result.Pinned, key)
		}
	}

	return entries, result, nil
}

// ValidateLock checks that all pinned keys in the LockFile exist in entries
// and that their values match. Returns a list of violation messages.
func ValidateLock(entries []Entry, lf *LockFile) []string {
	m := ToMap(entries)
	var violations []string
	for key, locked := range lf.Entries {
		val, ok := m[key]
		if !ok {
			violations = append(violations, fmt.Sprintf("pinned key %q is missing", key))
			continue
		}
		if val != locked.Value {
			violations = append(violations, fmt.Sprintf("pinned key %q has unexpected value", key))
		}
	}
	return violations
}
