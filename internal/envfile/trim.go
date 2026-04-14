package envfile

import (
	"strings"
)

// TrimOptions controls which entries are trimmed and how.
type TrimOptions struct {
	// TrimKeys removes leading/trailing whitespace from keys.
	TrimKeys bool
	// TrimValues removes leading/trailing whitespace from values.
	TrimValues bool
	// SkipSensitive prevents trimming of sensitive values.
	SkipSensitive bool
}

// TrimResult describes the outcome of a trim operation on a single entry.
type TrimResult struct {
	Key      string
	OldValue string
	NewValue string
	Changed  bool
	Skipped  bool
}

// Trim removes leading and trailing whitespace from keys and/or values
// according to the provided options. It returns the updated entries and
// a slice of TrimResult describing what changed.
func Trim(entries []Entry, opts TrimOptions) ([]Entry, []TrimResult) {
	results := make([]TrimResult, 0, len(entries))
	updated := make([]Entry, len(entries))

	for i, e := range entries {
		res := TrimResult{
			Key:      e.Key,
			OldValue: e.Value,
			NewValue: e.Value,
		}

		if opts.SkipSensitive && IsSensitive(e.Key) {
			res.Skipped = true
			updated[i] = e
			results = append(results, res)
			continue
		}

		newKey := e.Key
		if opts.TrimKeys {
			newKey = strings.TrimSpace(e.Key)
		}

		newVal := e.Value
		if opts.TrimValues {
			newVal = strings.TrimSpace(e.Value)
		}

		res.Key = newKey
		res.NewValue = newVal
		res.Changed = newKey != e.Key || newVal != e.Value

		updated[i] = Entry{Key: newKey, Value: newVal}
		results = append(results, res)
	}

	return updated, results
}

// TrimSummary returns a short human-readable summary of trim results.
func TrimSummary(results []TrimResult) string {
	changed, skipped := 0, 0
	for _, r := range results {
		if r.Changed {
			changed++
		}
		if r.Skipped {
			skipped++
		}
	}
	var sb strings.Builder
	sb.WriteString("trim: ")
	sb.WriteString(itoa(changed))
	sb.WriteString(" changed")
	if skipped > 0 {
		sb.WriteString(", ")
		sb.WriteString(itoa(skipped))
		sb.WriteString(" skipped")
	}
	return sb.String()
}

func itoa(n int) string {
	return strings.TrimSpace(strings.ReplaceAll(" "+string(rune('0'+n)), " ", ""))
}
