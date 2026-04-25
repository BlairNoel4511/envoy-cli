package envfile

// DefaultOptions controls how defaults are applied to entries.
type DefaultOptions struct {
	// Overwrite replaces existing values with the default.
	Overwrite bool
	// SkipSensitive prevents overwriting sensitive keys.
	SkipSensitive bool
}

// DefaultResult describes what happened to a single key during SetDefault.
type DefaultResult struct {
	Key     string
	Value   string
	Default string
	Status  string // "applied", "skipped", "unchanged", "not_found"
}

// DefaultSummary holds aggregate counts from a SetDefault run.
type DefaultSummary struct {
	Applied   int
	Skipped   int
	Unchanged int
	NotFound  int
}

// SetDefault applies a default value to a key if it is missing or empty.
// When opts.Overwrite is true it also replaces non-empty values.
func SetDefault(entries []Entry, key, defaultVal string, opts DefaultOptions) ([]Entry, DefaultResult) {
	for i, e := range entries {
		if e.Key != key {
			continue
		}
		if opts.SkipSensitive && IsSensitive(e.Key) {
			return entries, DefaultResult{Key: key, Value: e.Value, Default: defaultVal, Status: "skipped"}
		}
		if e.Value == defaultVal {
			return entries, DefaultResult{Key: key, Value: e.Value, Default: defaultVal, Status: "unchanged"}
		}
		if e.Value != "" && !opts.Overwrite {
			return entries, DefaultResult{Key: key, Value: e.Value, Default: defaultVal, Status: "skipped"}
		}
		entries[i].Value = defaultVal
		return entries, DefaultResult{Key: key, Value: defaultVal, Default: defaultVal, Status: "applied"}
	}
	// Key not found — append it.
	entries = append(entries, Entry{Key: key, Value: defaultVal})
	return entries, DefaultResult{Key: key, Value: defaultVal, Default: defaultVal, Status: "applied"}
}

// SetDefaults applies multiple key→default pairs in order.
func SetDefaults(entries []Entry, defaults map[string]string, opts DefaultOptions) ([]Entry, []DefaultResult, DefaultSummary) {
	var results []DefaultResult
	var sum DefaultSummary
	for key, val := range defaults {
		var res DefaultResult
		entries, res = SetDefault(entries, key, val, opts)
		results = append(results, res)
		switch res.Status {
		case "applied":
			sum.Applied++
		case "skipped":
			sum.Skipped++
		case "unchanged":
			sum.Unchanged++
		case "not_found":
			sum.NotFound++
		}
	}
	return entries, results, sum
}
