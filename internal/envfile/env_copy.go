package envfile

// CopyResult describes the outcome of copying a single key.
type CopyResult struct {
	SourceKey string
	DestKey   string
	Value     string
	Status    string // "copied", "skipped", "overwritten", "not_found"
}

// CopyOptions controls behaviour of Copy.
type CopyOptions struct {
	Overwrite      bool
	RedactSensitive bool
}

// Copy duplicates the value of srcKey into destKey within entries.
// Returns updated entries and a CopyResult.
func Copy(entries []Entry, srcKey, destKey string, opts CopyOptions) ([]Entry, CopyResult) {
	src, ok := Lookup(entries, srcKey)
	if !ok {
		return entries, CopyResult{SourceKey: srcKey, DestKey: destKey, Status: "not_found"}
	}

	value := src.Value
	if opts.RedactSensitive && IsSensitive(srcKey) {
		value = "***"
	}

	for i, e := range entries {
		if e.Key == destKey {
			if !opts.Overwrite {
				return entries, CopyResult{SourceKey: srcKey, DestKey: destKey, Value: value, Status: "skipped"}
			}
			entries[i].Value = src.Value
			return entries, CopyResult{SourceKey: srcKey, DestKey: destKey, Value: value, Status: "overwritten"}
		}
	}

	entries = append(entries, Entry{Key: destKey, Value: src.Value})
	return entries, CopyResult{SourceKey: srcKey, DestKey: destKey, Value: value, Status: "copied"}
}

// CopySummary returns a human-readable summary line.
func CopySummary(result CopyResult) string {
	switch result.Status {
	case "copied":
		return "copied " + result.SourceKey + " → " + result.DestKey
	case "overwritten":
		return "overwritten " + result.DestKey + " with value from " + result.SourceKey
	case "skipped":
		return "skipped: " + result.DestKey + " already exists"
	case "not_found":
		return "error: source key " + result.SourceKey + " not found"
	}
	return ""
}
