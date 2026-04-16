package envfile

import (
	"strings"
)

// SearchOptions controls how key/value search is performed.
type SearchOptions struct {
	CaseSensitive bool
	SearchKeys    bool
	SearchValues  bool
	Redact        bool
}

// SearchResult holds a matched entry and match context.
type SearchResult struct {
	Entry   Entry
	MatchedKey   bool
	MatchedValue bool
}

// SearchSummary holds aggregate counts.
type SearchSummary struct {
	Query   string
	Matched int
	Total   int
}

// Search scans entries for a query string in keys and/or values.
func Search(entries []Entry, query string, opts SearchOptions) []SearchResult {
	if !opts.SearchKeys && !opts.SearchValues {
		opts.SearchKeys = true
		opts.SearchValues = true
	}

	norm := query
	if !opts.CaseSensitive {
		norm = strings.ToLower(query)
	}

	var results []SearchResult
	for _, e := range entries {
		key := e.Key
		val := e.Value
		if !opts.CaseSensitive {
			key = strings.ToLower(e.Key)
			val = strings.ToLower(e.Value)
		}

		matchedKey := opts.SearchKeys && strings.Contains(key, norm)
		matchedVal := opts.SearchValues && strings.Contains(val, norm)

		if matchedKey || matchedVal {
			entry := e
			if opts.Redact && IsSensitive(e.Key) {
				entry.Value = "***"
			}
			results = append(results, SearchResult{
				Entry:        entry,
				MatchedKey:   matchedKey,
				MatchedValue: matchedVal,
			})
		}
	}
	return results
}

// SearchSummaryFor builds a summary from results.
func SearchSummaryFor(query string, results []SearchResult, total int) SearchSummary {
	return SearchSummary{
		Query:   query,
		Matched: len(results),
		Total:   total,
	}
}
