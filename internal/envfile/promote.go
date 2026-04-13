package envfile

import "fmt"

// PromoteOptions controls the behaviour of a promotion between profiles.
type PromoteOptions struct {
	// Overwrite allows existing keys in the target profile to be replaced.
	Overwrite bool
	// SkipSensitive prevents sensitive keys from being promoted.
	SkipSensitive bool
	// DryRun reports what would change without mutating any entries.
	DryRun bool
}

// PromoteResult describes what happened to a single key during promotion.
type PromoteResult struct {
	Key       string
	OldValue  string
	NewValue  string
	Action    string // "added", "overwritten", "skipped", "skipped_sensitive"
}

// PromoteSummary is the aggregate outcome of a Promote call.
type PromoteSummary struct {
	Results []PromoteResult
	DryRun  bool
}

// HasChanges returns true when at least one key was added or overwritten.
func (s PromoteSummary) HasChanges() bool {
	for _, r := range s.Results {
		if r.Action == "added" || r.Action == "overwritten" {
			return true
		}
	}
	return false
}

// Promote copies entries from src into dst according to opts.
// It returns the updated destination entries and a summary of all decisions.
func Promote(src, dst []Entry, opts PromoteOptions) ([]Entry, PromoteSummary) {
	dstMap := ToMap(dst)
	summary := PromoteSummary{DryRun: opts.DryRun}

	for _, e := range src {
		if opts.SkipSensitive && IsSensitive(e.Key) {
			summary.Results = append(summary.Results, PromoteResult{
				Key:    e.Key,
				Action: "skipped_sensitive",
			})
			continue
		}

		existing, exists := dstMap[e.Key]
		switch {
		case !exists:
			summary.Results = append(summary.Results, PromoteResult{
				Key:      e.Key,
				NewValue: e.Value,
				Action:   "added",
			})
			if !opts.DryRun {
				dstMap[e.Key] = e.Value
			}
		case exists && opts.Overwrite && existing != e.Value:
			summary.Results = append(summary.Results, PromoteResult{
				Key:      e.Key,
				OldValue: existing,
				NewValue: e.Value,
				Action:   "overwritten",
			})
			if !opts.DryRun {
				dstMap[e.Key] = e.Value
			}
		default:
			summary.Results = append(summary.Results, PromoteResult{
				Key:    e.Key,
				Action: "skipped",
			})
		}
	}

	return FromMap(dstMap), summary
}

// FormatPromoteSummary returns a human-readable summary of a promotion.
func FormatPromoteSummary(s PromoteSummary) string {
	if len(s.Results) == 0 {
		return "nothing to promote"
	}
	var added, overwritten, skipped int
	for _, r := range s.Results {
		switch r.Action {
		case "added":
			added++
		case "overwritten":
			overwritten++
		default:
			skipped++
		}
	}
	dryTag := ""
	if s.DryRun {
		dryTag = " (dry run)"
	}
	return fmt.Sprintf("promoted%s: %d added, %d overwritten, %d skipped",
		dryTag, added, overwritten, skipped)
}
