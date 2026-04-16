package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// ListOptions controls how entries are listed.
type ListOptions struct {
	SortKeys      bool
	RedactSecrets bool
	FilterPrefix  string
	ShowIndex     bool
}

// ListResult holds a single entry's display info.
type ListResult struct {
	Index int
	Key   string
	Value string
	Masked bool
}

// List returns display-ready results for the given entries.
func List(entries []Entry, opts ListOptions) []ListResult {
	filtered := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if opts.FilterPrefix != "" && !strings.HasPrefix(e.Key, opts.FilterPrefix) {
			continue
		}
		filtered = append(filtered, e)
	}

	if opts.SortKeys {
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].Key < filtered[j].Key
		})
	}

	results := make([]ListResult, 0, len(filtered))
	for i, e := range filtered {
		val := e.Value
		masked := false
		if opts.RedactSecrets && IsSensitive(e.Key) {
			val = "***"
			masked = true
		}
		results = append(results, ListResult{
			Index:  i + 1,
			Key:    e.Key,
			Value:  val,
			Masked: masked,
		})
	}
	return results
}

// ListSummary returns a short summary string.
func ListSummary(results []ListResult) string {
	masked := 0
	for _, r := range results {
		if r.Masked {
			masked++
		}
	}
	if masked > 0 {
		return fmt.Sprintf("%d entries (%d redacted)", len(results), masked)
	}
	return fmt.Sprintf("%d entries", len(results))
}
