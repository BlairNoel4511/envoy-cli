package envfile

// MergeOptions controls how two sets of entries are merged.
type MergeOptions struct {
	Overwrite      bool
	SkipSensitive  bool
	DryRun         bool
}

// MergeResult describes what happened to a single key during a merge.
type MergeResult struct {
	Key      string
	Value    string
	Action   string // "added", "updated", "skipped", "unchanged"
	Sensitive bool
}

// MergeSummary holds aggregate counts from a merge operation.
type MergeSummary struct {
	Added     int
	Updated   int
	Skipped   int
	Unchanged int
}

// MergeMany merges src entries into dst entries and returns the updated
// slice plus a per-key result list.
func MergeMany(dst, src []Entry, opts MergeOptions) ([]Entry, []MergeResult) {
	index := make(map[string]int, len(dst))
	for i, e := range dst {
		index[e.Key] = i
	}

	out := make([]Entry, len(dst))
	copy(out, dst)

	var results []MergeResult

	for _, s := range src {
		sensitive := IsSensitive(s.Key)
		if opts.SkipSensitive && sensitive {
			results = append(results, MergeResult{Key: s.Key, Value: s.Value, Action: "skipped", Sensitive: true})
			continue
		}

		if idx, exists := index[s.Key]; exists {
			if out[idx].Value == s.Value {
				results = append(results, MergeResult{Key: s.Key, Value: s.Value, Action: "unchanged", Sensitive: sensitive})
				continue
			}
			if !opts.Overwrite {
				results = append(results, MergeResult{Key: s.Key, Value: s.Value, Action: "skipped", Sensitive: sensitive})
				continue
			}
			if !opts.DryRun {
				out[idx].Value = s.Value
			}
			results = append(results, MergeResult{Key: s.Key, Value: s.Value, Action: "updated", Sensitive: sensitive})
		} else {
			if !opts.DryRun {
				out = append(out, s)
				index[s.Key] = len(out) - 1
			}
			results = append(results, MergeResult{Key: s.Key, Value: s.Value, Action: "added", Sensitive: sensitive})
		}
	}

	return out, results
}

// MergeManySum aggregates a result slice into a MergeSummary.
func MergeManySum(results []MergeResult) MergeSummary {
	var s MergeSummary
	for _, r := range results {
		switch r.Action {
		case "added":
			s.Added++
		case "updated":
			s.Updated++
		case "skipped":
			s.Skipped++
		case "unchanged":
			s.Unchanged++
		}
	}
	return s
}
