package envfile

import "time"

// TouchResult represents the outcome of touching (timestamping) a key.
type TouchResult struct {
	Key       string
	PrevTouch time.Time
	NewTouch  time.Time
	WasSet    bool
	Skipped   bool
	Reason    string
}

// TouchSummary holds aggregate counts for a touch operation.
type TouchSummary struct {
	Touched int
	Skipped int
}

// TouchOptions controls the behaviour of Touch.
type TouchOptions struct {
	// Overwrite allows updating the timestamp even if one is already set.
	Overwrite bool
	// Keys is the explicit list of keys to touch. If empty, all keys are touched.
	Keys []string
	// At overrides the timestamp; defaults to time.Now() when zero.
	At time.Time
}

// Touch records a "last touched" timestamp against env entries via a TouchStore.
// It returns per-key results and an aggregate summary.
func Touch(entries []Entry, store *TouchStore, opts TouchOptions) ([]TouchResult, TouchSummary) {
	now := opts.At
	if now.IsZero() {
		now = time.Now().UTC()
	}

	targetSet := make(map[string]struct{})
	for _, k := range opts.Keys {
		targetSet[k] = struct{}{}
	}

	var results []TouchResult
	var summary TouchSummary

	for _, e := range entries {
		if len(targetSet) > 0 {
			if _, ok := targetSet[e.Key]; !ok {
				continue
			}
		}

		prev, exists := store.Get(e.Key)
		if exists && !opts.Overwrite {
			results = append(results, TouchResult{
				Key:       e.Key,
				PrevTouch: prev,
				Skipped:   true,
				Reason:    "already touched",
			})
			summary.Skipped++
			continue
		}

		store.Set(e.Key, now)
		results = append(results, TouchResult{
			Key:      e.Key,
			PrevTouch: prev,
			NewTouch: now,
			WasSet:   true,
		})
		summary.Touched++
	}

	return results, summary
}
