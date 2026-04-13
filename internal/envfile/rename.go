package envfile

import "fmt"

// RenameResult describes the outcome of a single key rename operation.
type RenameResult struct {
	OldKey  string
	NewKey  string
	Value   string
	Skipped bool
	Reason  string
}

// RenameOptions controls the behaviour of Rename.
type RenameOptions struct {
	// Overwrite allows the new key to replace an existing entry.
	Overwrite bool
	// DryRun reports what would happen without modifying entries.
	DryRun bool
}

// Rename renames oldKey to newKey in entries, returning updated entries and a
// RenameResult that describes what happened. The original slice is never
// mutated; a new slice is always returned.
func Rename(entries []Entry, oldKey, newKey string, opts RenameOptions) ([]Entry, RenameResult, error) {
	if oldKey == "" || newKey == "" {
		return entries, RenameResult{}, fmt.Errorf("rename: key names must not be empty")
	}
	if oldKey == newKey {
		return entries, RenameResult{OldKey: oldKey, NewKey: newKey, Skipped: true, Reason: "old and new key are identical"}, nil
	}

	oldIdx := -1
	newIdx := -1
	for i, e := range entries {
		switch e.Key {
		case oldKey:
			oldIdx = i
		case newKey:
			newIdx = i
		}
	}

	if oldIdx == -1 {
		return entries, RenameResult{OldKey: oldKey, NewKey: newKey, Skipped: true, Reason: "key not found"}, nil
	}

	if newIdx != -1 && !opts.Overwrite {
		return entries, RenameResult{
			OldKey:  oldKey,
			NewKey:  newKey,
			Skipped: true,
			Reason:  "new key already exists (use Overwrite to replace)",
		}, nil
	}

	value := entries[oldIdx].Value
	result := RenameResult{OldKey: oldKey, NewKey: newKey, Value: value}

	if opts.DryRun {
		return entries, result, nil
	}

	// Build a new slice: replace oldKey with newKey; drop newIdx if overwriting.
	out := make([]Entry, 0, len(entries))
	for i, e := range entries {
		if i == newIdx {
			continue // will be overwritten by the renamed entry
		}
		if i == oldIdx {
			out = append(out, Entry{Key: newKey, Value: value})
			continue
		}
		out = append(out, e)
	}
	return out, result, nil
}
