package envfile

// MoveOptions controls behavior of the Move operation.
type MoveOptions struct {
	Overwrite bool
	DryRun    bool
}

// MoveResult describes the outcome of a single Move operation.
type MoveResult struct {
	From    string
	To      string
	Status  string // "moved", "skipped", "not_found", "conflict"
	Comment string
}

// MoveSummary holds aggregate counts for a Move operation.
type MoveSummary struct {
	Moved    int
	Skipped  int
	NotFound int
	Conflict int
}

// Move renames srcKey to dstKey in entries, removing the source.
// It returns updated entries, a result, and a summary.
func Move(entries []Entry, srcKey, dstKey string, opts MoveOptions) ([]Entry, MoveResult, MoveSummary) {
	result := MoveResult{From: srcKey, To: dstKey}

	srcIdx := -1
	dstIdx := -1
	for i, e := range entries {
		if e.Key == srcKey {
			srcIdx = i
		}
		if e.Key == dstKey {
			dstIdx = i
		}
	}

	if srcIdx == -1 {
		result.Status = "not_found"
		result.Comment = "source key does not exist"
		return entries, result, MoveSummary{NotFound: 1}
	}

	if dstIdx != -1 && !opts.Overwrite {
		result.Status = "conflict"
		result.Comment = "destination key already exists"
		return entries, result, MoveSummary{Conflict: 1}
	}

	result.Status = "moved"
	if opts.DryRun {
		result.Comment = "dry-run"
		return entries, result, MoveSummary{Moved: 1}
	}

	srcEntry := entries[srcIdx]
	// Remove source
	updated := make([]Entry, 0, len(entries))
	for i, e := range entries {
		if i == srcIdx {
			continue
		}
		if i == dstIdx {
			continue
		}
		updated = append(updated, e)
	}
	updated = append(updated, Entry{Key: dstKey, Value: srcEntry.Value})
	return updated, result, MoveSummary{Moved: 1}
}
