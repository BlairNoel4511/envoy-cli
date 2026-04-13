package envfile

import "fmt"

// CloneOptions controls how a clone operation behaves.
type CloneOptions struct {
	// Overwrite replaces existing keys in the destination.
	Overwrite bool
	// SkipSensitive prevents sensitive keys from being cloned.
	SkipSensitive bool
	// Prefix filters source keys to only those with the given prefix.
	Prefix string
}

// CloneResult holds the outcome of a Clone operation.
type CloneResult struct {
	Cloned  []string
	Skipped []string
	Errors  []string
}

// hasPrefix reports whether key starts with prefix.
func hasPrefix(key, prefix string) bool {
	return len(key) >= len(prefix) && key[:len(prefix)] == prefix
}

// Clone copies entries from src into dst according to the given options.
// It returns a CloneResult summarising what happened.
func Clone(src, dst []Entry, opts CloneOptions) ([]Entry, CloneResult) {
	result := CloneResult{}

	dstMap := ToMap(dst)

	out := make([]Entry, len(dst))
	copy(out, dst)

	for _, e := range src {
		if opts.Prefix != "" && !hasPrefix(e.Key, opts.Prefix) {
			continue
		}

		if opts.SkipSensitive && IsSensitive(e.Key) {
			result.Skipped = append(result.Skipped, e.Key)
			continue
		}

		if _, exists := dstMap[e.Key]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, e.Key)
			continue
		}

		if _, exists := dstMap[e.Key]; exists && opts.Overwrite {
			for i, d := range out {
				if d.Key == e.Key {
					out[i].Value = e.Value
					break
				}
			}
		} else {
			out = append(out, e)
		}

		dstMap[e.Key] = e.Value
		result.Cloned = append(result.Cloned, e.Key)
	}

	return out, result
}

// FormatCloneResult returns a human-readable summary of a CloneResult.
func FormatCloneResult(r CloneResult) string {
	var s string
	s += fmt.Sprintf("Cloned:  %d key(s)\n", len(r.Cloned))
	s += fmt.Sprintf("Skipped: %d key(s)\n", len(r.Skipped))
	if len(r.Errors) > 0 {
		s += fmt.Sprintf("Errors:  %d\n", len(r.Errors))
		for _, e := range r.Errors {
			s += fmt.Sprintf("  - %s\n", e)
		}
	}
	return s
}
