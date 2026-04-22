package envfile

// DiffApplyOptions controls how a diff is applied to a set of entries.
type DiffApplyOptions struct {
	Overwrite bool
	SkipSensitive bool
	DryRun bool
}

// DiffApplyResult captures the outcome of applying a single diff change.
type DiffApplyResult struct {
	Key string
	Action string // "added", "updated", "removed", "skipped"
	OldValue string
	NewValue string
	Sensitive bool
}

// DiffApplySummary holds aggregate counts from a DiffApply operation.
type DiffApplySummary struct {
	Added int
	Updated int
	Removed int
	Skipped int
}

// ApplyDiff applies a set of DiffEntry changes onto a base slice of Entry.
// It returns the modified entries and a list of per-key results.
func ApplyDiff(base []Entry, changes []DiffEntry, opts DiffApplyOptions) ([]Entry, []DiffApplyResult) {
	entries := make([]Entry, len(base))
	copy(entries, base)

	var results []DiffApplyResult

	for _, change := range changes {
		sensitive := IsSensitive(change.Key)
		if opts.SkipSensitive && sensitive {
			results = append(results, DiffApplyResult{
				Key: change.Key, Action: "skipped", Sensitive: true,
			})
			continue
		}

		switch change.Status {
		case "added":
			if !opts.DryRun {
				entries = append(entries, Entry{Key: change.Key, Value: change.NewValue})
			}
			results = append(results, DiffApplyResult{Key: change.Key, Action: "added", NewValue: change.NewValue, Sensitive: sensitive})

		case "removed":
			if !opts.DryRun {
				entries = removeKey(entries, change.Key)
			}
			results = append(results, DiffApplyResult{Key: change.Key, Action: "removed", OldValue: change.OldValue, Sensitive: sensitive})

		case "changed":
			idx := findKeyIndex(entries, change.Key)
			if idx < 0 {
				results = append(results, DiffApplyResult{Key: change.Key, Action: "skipped", Sensitive: sensitive})
				continue
			}
			if !opts.Overwrite {
				results = append(results, DiffApplyResult{Key: change.Key, Action: "skipped", OldValue: change.OldValue, Sensitive: sensitive})
				continue
			}
			if !opts.DryRun {
				entries[idx].Value = change.NewValue
			}
			results = append(results, DiffApplyResult{Key: change.Key, Action: "updated", OldValue: change.OldValue, NewValue: change.NewValue, Sensitive: sensitive})
		}
	}

	return entries, results
}

// SummarizeDiffApply aggregates a slice of DiffApplyResult into a DiffApplySummary.
func SummarizeDiffApply(results []DiffApplyResult) DiffApplySummary {
	var s DiffApplySummary
	for _, r := range results {
		switch r.Action {
		case "added":
			s.Added++
		case "updated":
			s.Updated++
		case "removed":
			s.Removed++
		case "skipped":
			s.Skipped++
		}
	}
	return s
}

func removeKey(entries []Entry, key string) []Entry {
	out := entries[:0]
	for _, e := range entries {
		if e.Key != key {
			out = append(out, e)
		}
	}
	return out
}

func findKeyIndex(entries []Entry, key string) int {
	for i, e := range entries {
		if e.Key == key {
			return i
		}
	}
	return -1
}
