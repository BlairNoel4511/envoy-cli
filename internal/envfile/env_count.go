package envfile

import (
	"strings"
)

// CountOptions controls how entries are counted.
type CountOptions struct {
	Prefix    string
	Sensitive bool
	NonEmpty  bool
}

// CountResult holds the result of a Count operation.
type CountResult struct {
	Total     int
	Sensitive int
	Empty     int
	NonEmpty  int
	Filtered  int
}

// Count counts entries in the given slice according to the provided options.
// If a prefix is set, only entries whose keys start with that prefix are counted.
// If SensitiveOnly is true, only sensitive keys are counted.
// If NonEmpty is true, only entries with non-empty values are counted.
func Count(entries []Entry, opts CountOptions) CountResult {
	var result CountResult

	for _, e := range entries {
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}

		result.Total++

		if IsSensitive(e.Key) {
			result.Sensitive++
		}

		if e.Value == "" {
			result.Empty++
		} else {
			result.NonEmpty++
		}
	}

	if opts.Sensitive {
		result.Filtered = result.Sensitive
	} else if opts.NonEmpty {
		result.Filtered = result.NonEmpty
	} else {
		result.Filtered = result.Total
	}

	return result
}

// CountSummary returns a human-readable summary string for a CountResult.
func CountSummary(r CountResult) string {
	parts := []string{}
	parts = append(parts, itoa(r.Total)+" total")
	parts = append(parts, itoa(r.Sensitive)+" sensitive")
	parts = append(parts, itoa(r.Empty)+" empty")
	parts = append(parts, itoa(r.NonEmpty)+" non-empty")

	out := ""
	for i, p := range parts {
		if i > 0 {
			out += ", "
		}
		out += p
	}
	return out
}
