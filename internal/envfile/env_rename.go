package envfile

// RenameResult describes the outcome of a single rename operation.
type RenameResult struct {
	OldKey  string
	NewKey  string
	Status  string // "renamed", "not_found", "conflict", "skipped"
	Comment string
}

// RenameManyOptions controls batch rename behaviour.
type RenameManyOptions struct {
	Overwrite bool
	DryRun    bool
}

// RenameMany renames multiple keys in entries according to the provided
// oldKey→newKey mapping. It returns one RenameResult per mapping entry.
func RenameMany(entries []Entry, mapping map[string]string, opts RenameManyOptions) ([]Entry, []RenameResult) {
	results := make([]RenameResult, 0, len(mapping))
	working := make([]Entry, len(entries))
	copy(working, entries)

	for oldKey, newKey := range mapping {
		var found bool
		for _, e := range working {
			if e.Key == oldKey {
				found = true
				break
			}
		}
		if !found {
			results = append(results, RenameResult{OldKey: oldKey, NewKey: newKey, Status: "not_found"})
			continue
		}

		_, exists := Lookup(working, newKey)
		if exists && !opts.Overwrite {
			results = append(results, RenameResult{OldKey: oldKey, NewKey: newKey, Status: "conflict", Comment: "target key already exists"})
			continue
		}

		if opts.DryRun {
			results = append(results, RenameResult{OldKey: oldKey, NewKey: newKey, Status: "skipped", Comment: "dry run"})
			continue
		}

		updated := make([]Entry, 0, len(working))
		for _, e := range working {
			if e.Key == newKey && exists {
				continue // remove old occupant
			}
			if e.Key == oldKey {
				e.Key = newKey
			}
			updated = append(updated, e)
		}
		working = updated
		results = append(results, RenameResult{OldKey: oldKey, NewKey: newKey, Status: "renamed"})
	}

	return working, results
}

// RenameSummary returns counts of each status from a slice of RenameResult.
func RenameSummary(results []RenameResult) map[string]int {
	m := map[string]int{}
	for _, r := range results {
		m[r.Status]++
	}
	return m
}
