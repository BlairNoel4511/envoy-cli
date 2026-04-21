package envfile

import (
	"fmt"
	"sort"
)

// ProtectOptions configures protection behavior.
type ProtectOptions struct {
	Overwrite bool
}

// ProtectResult describes the outcome of protecting a single key.
type ProtectResult struct {
	Key       string
	Protected bool
	Already   bool
	NotFound  bool
}

// ProtectSummary holds aggregate counts for a protect operation.
type ProtectSummary struct {
	Protected int
	Already   int
	NotFound  int
}

// Protect marks one or more keys as protected in the given entries.
// Protected keys cannot be overwritten by sync, import, or merge unless
// explicitly forced. Returns per-key results and a summary.
func Protect(entries []Entry, keys []string, store *ProtectStore, opts ProtectOptions) ([]ProtectResult, ProtectSummary) {
	keySet := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		keySet[e.Key] = struct{}{}
	}

	var results []ProtectResult
	var summary ProtectSummary

	for _, k := range keys {
		if _, exists := keySet[k]; !exists {
			results = append(results, ProtectResult{Key: k, NotFound: true})
			summary.NotFound++
			continue
		}
		if store.IsProtected(k) && !opts.Overwrite {
			results = append(results, ProtectResult{Key: k, Already: true})
			summary.Already++
			continue
		}
		store.Add(k)
		results = append(results, ProtectResult{Key: k, Protected: true})
		summary.Protected++
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results, summary
}

// ProtectSummaryText returns a one-line human-readable summary.
func ProtectSummaryText(s ProtectSummary) string {
	return fmt.Sprintf("protected %d, already protected %d, not found %d",
		s.Protected, s.Already, s.NotFound)
}
