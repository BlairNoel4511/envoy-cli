package envfile

// UnsetResult describes the outcome of unsetting a key.
type UnsetResult struct {
	Key     string
	Removed bool
	Missing bool
}

// UnsetSummary holds aggregate counts for an unset operation.
type UnsetSummary struct {
	Removed int
	Missing int
}

// Unset removes a single key from entries, returning the updated slice and result.
func Unset(entries []Entry, key string) ([]Entry, UnsetResult) {
	result := UnsetResult{Key: key}
	updated := make([]Entry, 0, len(entries))
	found := false
	for _, e := range entries {
		if e.Key == key {
			found = true
			continue
		}
		updated = append(updated, e)
	}
	if found {
		result.Removed = true
	} else {
		result.Missing = true
	}
	return updated, result
}

// UnsetMany removes multiple keys from entries.
func UnsetMany(entries []Entry, keys []string) ([]Entry, []UnsetResult, UnsetSummary) {
	results := make([]UnsetResult, 0, len(keys))
	summary := UnsetSummary{}
	current := entries
	for _, k := range keys {
		var res UnsetResult
		current, res = Unset(current, k)
		results = append(results, res)
		if res.Removed {
			summary.Removed++
		} else {
			summary.Missing++
		}
	}
	return current, results, summary
}
