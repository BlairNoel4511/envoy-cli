package envfile

// DuplicateResult holds the outcome of a duplicate operation for a single key.
type DuplicateResult struct {
	SourceKey string
	DestKey   string
	Status    string // "duplicated", "skipped", "source_not_found", "unchanged"
	Value     string
	Sensitive bool
}

// DuplicateSummary holds aggregate counts for a duplicate operation.
type DuplicateSummary struct {
	Duplicated    int
	Skipped       int
	NotFound      int
	Unchanged     int
}

// DuplicateOptions controls the behaviour of Duplicate.
type DuplicateOptions struct {
	Overwrite      bool
	SkipSensitive  bool
}

// Duplicate copies the value of srcKey into destKey within entries.
// It returns the updated entries, a result record, and a summary.
func Duplicate(entries []Entry, srcKey, destKey string, opts DuplicateOptions) ([]Entry, DuplicateResult, DuplicateSummary) {
	result := DuplicateResult{SourceKey: srcKey, DestKey: destKey}

	// Find source
	srcVal, srcFound := Lookup(entries, srcKey)
	if !srcFound {
		result.Status = "source_not_found"
		return entries, result, DuplicateSummary{NotFound: 1}
	}

	isSensitive := IsSensitive(srcKey)
	result.Sensitive = isSensitive

	if opts.SkipSensitive && isSensitive {
		result.Status = "skipped"
		return entries, result, DuplicateSummary{Skipped: 1}
	}

	// Check if dest already exists
	existingVal, destFound := Lookup(entries, destKey)
	if destFound {
		if existingVal == srcVal {
			result.Status = "unchanged"
			result.Value = srcVal
			return entries, result, DuplicateSummary{Unchanged: 1}
		}
		if !opts.Overwrite {
			result.Status = "skipped"
			return entries, result, DuplicateSummary{Skipped: 1}
		}
		// Overwrite existing dest
		for i, e := range entries {
			if e.Key == destKey {
				entries[i].Value = srcVal
				break
			}
		}
		result.Status = "duplicated"
		result.Value = srcVal
		return entries, result, DuplicateSummary{Duplicated: 1}
	}

	// Append new entry
	entries = append(entries, Entry{Key: destKey, Value: srcVal})
	result.Status = "duplicated"
	result.Value = srcVal
	return entries, result, DuplicateSummary{Duplicated: 1}
}
