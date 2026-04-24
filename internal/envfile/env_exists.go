package envfile

// ExistsResult holds the result of checking whether a key exists.
type ExistsResult struct {
	Key     string
	Exists  bool
	Value   string
	Masked  bool
}

// ExistsOptions controls behaviour of the Exists check.
type ExistsOptions struct {
	RedactSensitive bool
}

// Exists checks whether a single key is present in the given entries.
func Exists(entries []Entry, key string, opts ExistsOptions) ExistsResult {
	for _, e := range entries {
		if e.Key == key {
			val := e.Value
			masked := false
			if opts.RedactSensitive && IsSensitive(e.Key) {
				val = "***"
				masked = true
			}
			return ExistsResult{Key: key, Exists: true, Value: val, Masked: masked}
		}
	}
	return ExistsResult{Key: key, Exists: false}
}

// ExistsMany checks whether each of the given keys is present in entries.
func ExistsMany(entries []Entry, keys []string, opts ExistsOptions) []ExistsResult {
	results := make([]ExistsResult, 0, len(keys))
	for _, k := range keys {
		results = append(results, Exists(entries, k, opts))
	}
	return results
}

// ExistsSummary returns counts of found and missing keys.
type ExistsSummary struct {
	Found   int
	Missing int
}

// SummarizeExists builds an ExistsSummary from a slice of results.
func SummarizeExists(results []ExistsResult) ExistsSummary {
	s := ExistsSummary{}
	for _, r := range results {
		if r.Exists {
			s.Found++
		} else {
			s.Missing++
		}
	}
	return s
}
