package envfile

import (
	"strings"
)

// FilterOptions controls how entries are filtered.
type FilterOptions struct {
	// Prefix filters entries whose keys start with the given prefix.
	Prefix string
	// SensitiveOnly returns only entries considered sensitive.
	SensitiveOnly bool
	// Keys is an explicit allowlist of keys to include. If empty, all keys pass.
	Keys []string
}

// Filter returns a subset of entries based on the provided FilterOptions.
func Filter(entries []Entry, opts FilterOptions) []Entry {
	allowlist := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		allowlist[k] = true
	}

	var result []Entry
	for _, e := range entries {
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if opts.SensitiveOnly && !IsSensitive(e.Key) {
			continue
		}
		if len(allowlist) > 0 && !allowlist[e.Key] {
			continue
		}
		result = append(result, e)
	}
	return result
}

// FilterByPrefix is a convenience wrapper around Filter for prefix-only filtering.
func FilterByPrefix(entries []Entry, prefix string) []Entry {
	return Filter(entries, FilterOptions{Prefix: prefix})
}

// FilterSensitive returns only entries whose keys are considered sensitive.
func FilterSensitive(entries []Entry) []Entry {
	return Filter(entries, FilterOptions{SensitiveOnly: true})
}
