package envfile

// DeleteResult represents the outcome of a delete operation on a single key.
type DeleteResult struct {
	Key     string
	Deleted bool
	Reason  string
}

// DeleteSummary holds aggregate counts for a batch delete.
type DeleteSummary struct {
	Deleted int
	Skipped int
}

// Delete removes a single key from entries. Returns updated entries and a result.
func Delete(entries []Entry, key string) ([]Entry, DeleteResult) {
	for i, e := range entries {
		if e.Key == key {
			updated := append(entries[:i:i], entries[i+1:]...)
			return updated, DeleteResult{Key: key, Deleted: true, Reason: "deleted"}
		}
	}
	return entries, DeleteResult{Key: key, Deleted: false, Reason: "key not found"}
}

// DeleteMany removes multiple keys from entries. Returns updated entries and per-key results.
func DeleteMany(entries []Entry, keys []string) ([]Entry, []DeleteResult, DeleteSummary) {
	results := make([]DeleteResult, 0, len(keys))
	summary := DeleteSummary{}

	current := entries
	for _, key := range keys {
		var res DeleteResult
		current, res = Delete(current, key)
		results = append(results, res)
		if res.Deleted {
			summary.Deleted++
		} else {
			summary.Skipped++
		}
	}
	return current, results, summary
}
