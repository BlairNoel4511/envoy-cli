package envfile

import "time"

// FreezeResult describes the outcome of freezing a single entry.
type FreezeResult struct {
	Key       string
	Frozen    bool
	Already   bool
	NotFound  bool
	Timestamp time.Time
}

// FreezeSummary holds aggregate counts for a freeze operation.
type FreezeSummary struct {
	Frozen   int
	Already  int
	NotFound int
}

// FreezeOptions controls the behaviour of Freeze.
type FreezeOptions struct {
	// Force re-freezes keys that are already frozen, resetting their timestamp.
	Force bool
}

// Freeze marks the given keys as frozen in the provided FreezeStore.
// Frozen keys cannot be modified by Set, Rotate, or Patch operations.
func Freeze(entries []Entry, store *FreezeStore, keys []string, opts FreezeOptions) ([]FreezeResult, FreezeSummary) {
	keySet := make(map[string]bool, len(entries))
	for _, e := range entries {
		keySet[e.Key] = true
	}

	var results []FreezeResult
	var sum FreezeSummary
	now := time.Now().UTC()

	for _, k := range keys {
		if !keySet[k] {
			results = append(results, FreezeResult{Key: k, NotFound: true})
			sum.NotFound++
			continue
		}
		if store.IsFrozen(k) && !opts.Force {
			results = append(results, FreezeResult{Key: k, Already: true})
			sum.Already++
			continue
		}
		store.Freeze(k, now)
		results = append(results, FreezeResult{Key: k, Frozen: true, Timestamp: now})
		sum.Frozen++
	}

	return results, sum
}

// FreezeSummaryText returns a human-readable one-line summary.
func FreezeSummaryText(s FreezeSummary) string {
	return "frozen: " + itoa(s.Frozen) +
		", already frozen: " + itoa(s.Already) +
		", not found: " + itoa(s.NotFound)
}
