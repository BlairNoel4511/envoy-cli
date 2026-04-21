package envfile

import (
	"sort"
	"strings"
)

// SortOrder defines the direction of sorting.
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// SortOptions controls how entries are sorted.
type SortOptions struct {
	Order     SortOrder
	ByValue   bool
	IgnoreCase bool
}

// SortResult holds the outcome of a sort operation.
type SortResult struct {
	Entries  []Entry
	Original []Entry
	Reordered int
}

// Sort returns a sorted copy of entries according to the given options.
// By default entries are sorted by key in ascending order.
func Sort(entries []Entry, opts SortOptions) SortResult {
	if opts.Order == "" {
		opts.Order = SortAsc
	}

	sorted := make([]Entry, len(entries))
	copy(sorted, entries)

	sort.SliceStable(sorted, func(i, j int) bool {
		var a, b string
		if opts.ByValue {
			a, b = sorted[i].Value, sorted[j].Value
		} else {
			a, b = sorted[i].Key, sorted[j].Key
		}
		if opts.IgnoreCase {
			a = strings.ToLower(a)
			b = strings.ToLower(b)
		}
		if opts.Order == SortDesc {
			return a > b
		}
		return a < b
	})

	reordered := 0
	for i := range entries {
		if i < len(sorted) && sorted[i].Key != entries[i].Key {
			reordered++
		}
	}

	return SortResult{
		Entries:   sorted,
		Original:  entries,
		Reordered: reordered,
	}
}

// SortSummary returns a human-readable summary of the sort result.
func SortSummary(r SortResult) string {
	if r.Reordered == 0 {
		return "already sorted, no changes"
	}
	return itoa(r.Reordered) + " of " + itoa(len(r.Entries)) + " entries reordered"
}
