package envfile

import "strings"

// ExtractOptions controls how entries are extracted into a new set.
type ExtractOptions struct {
	Keys        []string
	Prefix      string
	SensitiveOnly bool
	StripPrefix bool
}

// ExtractResult holds the outcome of a single extraction.
type ExtractResult struct {
	Key         string
	ExtractedAs string
	Value       string
	Sensitive   bool
	Skipped     bool
	Reason      string
}

// ExtractSummary holds aggregate counts for an extraction run.
type ExtractSummary struct {
	Extracted int
	Skipped   int
}

// Extract pulls a subset of entries based on the provided options.
// It returns the extracted entries and a per-key result slice.
func Extract(entries []Entry, opts ExtractOptions) ([]Entry, []ExtractResult, ExtractSummary) {
	allowset := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		allowset[k] = true
	}

	var out []Entry
	var results []ExtractResult
	var summary ExtractSummary

	for _, e := range entries {
		if opts.SensitiveOnly && !IsSensitive(e.Key) {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if len(allowset) > 0 && !allowset[e.Key] {
			results = append(results, ExtractResult{
				Key:     e.Key,
				Skipped: true,
				Reason:  "not in allowlist",
			})
			summary.Skipped++
			continue
		}

		extractedKey := e.Key
		if opts.StripPrefix && opts.Prefix != "" {
			extractedKey = strings.TrimPrefix(e.Key, opts.Prefix)
		}

		out = append(out, Entry{Key: extractedKey, Value: e.Value, Comment: e.Comment})
		results = append(results, ExtractResult{
			Key:         e.Key,
			ExtractedAs: extractedKey,
			Value:       e.Value,
			Sensitive:   IsSensitive(e.Key),
		})
		summary.Extracted++
	}

	return out, results, summary
}
