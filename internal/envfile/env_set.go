package envfile

import "fmt"

// SetOptions controls how individual key-value pairs are set.
type SetOptions struct {
	Overwrite bool
	DryRun    bool
}

// SetResult describes the outcome of a single Set operation.
type SetResult struct {
	Key     string
	Value   string
	OldValue string
	Action  string // "added", "updated", "skipped", "unchanged"
}

// Set applies a single key-value assignment to a slice of entries.
// It returns the updated entries and a SetResult describing what happened.
func Set(entries []Entry, key, value string, opts SetOptions) ([]Entry, SetResult) {
	result := SetResult{Key: key, Value: value}

	for i, e := range entries {
		if e.Key != key {
			continue
		}
		if e.Value == value {
			result.Action = "unchanged"
			result.OldValue = e.Value
			return entries, result
		}
		if !opts.Overwrite {
			result.Action = "skipped"
			result.OldValue = e.Value
			return entries, result
		}
		result.OldValue = e.Value
		result.Action = "updated"
		if !opts.DryRun {
			entries[i].Value = value
		}
		return entries, result
	}

	result.Action = "added"
	if !opts.DryRun {
		entries = append(entries, Entry{Key: key, Value: value})
	}
	return entries, result
}

// SetMany applies multiple key-value pairs using Set and returns all results.
func SetMany(entries []Entry, pairs map[string]string, opts SetOptions) ([]Entry, []SetResult) {
	results := make([]SetResult, 0, len(pairs))
	for k, v := range pairs {
		var r SetResult
		entries, r = Set(entries, k, v, opts)
		results = append(results, r)
	}
	return entries, results
}

// SetSummary returns a human-readable summary string for a slice of SetResults.
func SetSummary(results []SetResult) string {
	added, updated, skipped, unchanged := 0, 0, 0, 0
	for _, r := range results {
		switch r.Action {
		case "added":
			added++
		case "updated":
			updated++
		case "skipped":
			skipped++
		case "unchanged":
			unchanged++
		}
	}
	return fmt.Sprintf("added=%d updated=%d skipped=%d unchanged=%d", added, updated, skipped, unchanged)
}
