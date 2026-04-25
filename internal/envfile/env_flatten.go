package envfile

import (
	"strings"
)

// FlattenOptions controls how nested/prefixed keys are flattened.
type FlattenOptions struct {
	// Separator is the delimiter used to join prefix segments (default: "_").
	Separator string
	// StripPrefix removes the leading prefix from keys after flattening.
	StripPrefix bool
	// Uppercase normalizes all resulting keys to uppercase.
	Uppercase bool
	// Redact masks sensitive values in results.
	Redact bool
}

// FlattenResult holds the outcome for a single entry after flattening.
type FlattenResult struct {
	OriginalKey string
	NewKey      string
	Value       string
	Changed     bool
	Sensitive   bool
}

// FlattenSummary holds aggregate counts for a flatten operation.
type FlattenSummary struct {
	Total   int
	Changed int
	Skipped int
}

// Flatten normalizes keys by collapsing repeated separators, optionally
// stripping a common prefix and uppercasing all keys.
func Flatten(entries []Entry, opts FlattenOptions) ([]Entry, []FlattenResult, FlattenSummary) {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	out := make([]Entry, 0, len(entries))
	results := make([]FlattenResult, 0, len(entries))
	var summary FlattenSummary

	for _, e := range entries {
		summary.Total++
		sensitive := IsSensitive(e.Key)

		newKey := collapseRepeated(e.Key, opts.Separator)
		if opts.StripPrefix {
			newKey = strings.TrimPrefix(newKey, opts.Separator)
		}
		if opts.Uppercase {
			newKey = strings.ToUpper(newKey)
		}

		changed := newKey != e.Key
		if changed {
			summary.Changed++
		} else {
			summary.Skipped++
		}

		displayVal := e.Value
		if opts.Redact && sensitive {
			displayVal = "***"
		}

		results = append(results, FlattenResult{
			OriginalKey: e.Key,
			NewKey:      newKey,
			Value:       displayVal,
			Changed:     changed,
			Sensitive:   sensitive,
		})
		out = append(out, Entry{Key: newKey, Value: e.Value, Comment: e.Comment})
	}

	return out, results, summary
}

// collapseRepeated replaces runs of the separator character with a single instance.
func collapseRepeated(key, sep string) string {
	if sep == "" {
		return key
	}
	for strings.Contains(key, sep+sep) {
		key = strings.ReplaceAll(key, sep+sep, sep)
	}
	return strings.Trim(key, sep)
}
