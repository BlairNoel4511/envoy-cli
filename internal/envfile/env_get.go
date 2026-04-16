package envfile

import "fmt"

// GetOptions controls how Get behaves.
type GetOptions struct {
	Redact  bool
	Default string
}

// GetResult holds the outcome of a single Get operation.
type GetResult struct {
	Key     string
	Value   string
	Found   bool
	Redacted bool
}

// Get retrieves a single key from entries.
func Get(entries []Entry, key string, opts GetOptions) GetResult {
	for _, e := range entries {
		if e.Key == key {
			v := e.Value
			redacted := false
			if opts.Redact && IsSensitive(key) {
				v = "***"
				redacted = true
			}
			return GetResult{Key: key, Value: v, Found: true, Redacted: redacted}
		}
	}
	return GetResult{Key: key, Value: opts.Default, Found: false}
}

// GetMany retrieves multiple keys from entries.
func GetMany(entries []Entry, keys []string, opts GetOptions) []GetResult {
	results := make([]GetResult, 0, len(keys))
	for _, k := range keys {
		results = append(results, Get(entries, k, opts))
	}
	return results
}

// GetSummary returns a human-readable summary line.
func GetSummary(results []GetResult) string {
	found, missing := 0, 0
	for _, r := range results {
		if r.Found {
			found++
		} else {
			missing++
		}
	}
	return fmt.Sprintf("%d found, %d missing", found, missing)
}
