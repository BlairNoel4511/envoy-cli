package envfile

import "fmt"

// RotateOptions controls how key rotation behaves.
type RotateOptions struct {
	DryRun      bool
	SkipMissing bool
	AuditLog    *AuditLog
}

// RotateResult captures the outcome of a single key rotation.
type RotateResult struct {
	Key      string
	OldValue string
	NewValue string
	Skipped  bool
	Reason   string
}

// RotateSummary holds the aggregated results of a rotation run.
type RotateSummary struct {
	Results  []RotateResult
	Rotated  int
	Skipped  int
}

// Rotate replaces the values of the specified keys in entries with the
// corresponding values from newValues. Keys absent from entries are skipped
// when SkipMissing is true, otherwise an error is returned.
func Rotate(entries []Entry, newValues map[string]string, opts RotateOptions) ([]Entry, RotateSummary, error) {
	index := make(map[string]int, len(entries))
	for i, e := range entries {
		index[e.Key] = i
	}

	result := make([]Entry, len(entries))
	copy(result, entries)

	var summary RotateSummary

	for key, newVal := range newValues {
		i, found := index[key]
		if !found {
			if opts.SkipMissing {
				summary.Results = append(summary.Results, RotateResult{
					Key:     key,
					Skipped: true,
					Reason:  "key not found",
				})
				summary.Skipped++
				continue
			}
			return nil, RotateSummary{}, fmt.Errorf("rotate: key %q not found in entries", key)
		}

		oldVal := result[i].Value
		if oldVal == newVal {
			summary.Results = append(summary.Results, RotateResult{
				Key:     key,
				Skipped: true,
				Reason:  "value unchanged",
			})
			summary.Skipped++
			continue
		}

		if !opts.DryRun {
			result[i].Value = newVal
			if opts.AuditLog != nil {
				opts.AuditLog.Record("rotate", key, oldVal, newVal)
			}
		}

		summary.Results = append(summary.Results, RotateResult{
			Key:      key,
			OldValue: oldVal,
			NewValue: newVal,
		})
		summary.Rotated++
	}

	return result, summary, nil
}
