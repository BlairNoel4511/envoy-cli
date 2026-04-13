package envfile

import (
	"fmt"
	"strings"
)

// ResolveOption controls resolution behaviour.
type ResolveOption struct {
	AllowMissing bool   // if true, unresolved vars are left as-is
	Prefix       string // optional prefix filter; only resolve keys with this prefix
}

// ResolveResult holds the outcome of a resolve operation for one entry.
type ResolveResult struct {
	Key        string
	Original   string
	Resolved   string
	Changed    bool
	Unresolved []string // variable names that could not be resolved
}

// Resolve performs variable interpolation on entry values using the provided
// lookup map. References of the form $KEY or ${KEY} are replaced with their
// values from the map. Entries not matching the optional prefix are skipped.
func Resolve(entries []Entry, lookup map[string]string, opt ResolveOption) []ResolveResult {
	results := make([]ResolveResult, 0, len(entries))
	for _, e := range entries {
		if opt.Prefix != "" && !strings.HasPrefix(e.Key, opt.Prefix) {
			continue
		}
		resolved, missing := interpolate(e.Value, lookup)
		if !opt.AllowMissing && len(missing) > 0 {
			resolved = e.Value // leave original on error
		}
		results = append(results, ResolveResult{
			Key:        e.Key,
			Original:   e.Value,
			Resolved:   resolved,
			Changed:    resolved != e.Value,
			Unresolved: missing,
		})
	}
	return results
}

// ApplyResolved returns a new entry slice with resolved values applied.
func ApplyResolved(entries []Entry, results []ResolveResult) []Entry {
	patch := make(map[string]string, len(results))
	for _, r := range results {
		if r.Changed {
			patch[r.Key] = r.Resolved
		}
	}
	out := make([]Entry, len(entries))
	for i, e := range entries {
		if v, ok := patch[e.Key]; ok {
			out[i] = Entry{Key: e.Key, Value: v}
		} else {
			out[i] = e
		}
	}
	return out
}

// interpolate replaces $KEY / ${KEY} references in s using lookup.
func interpolate(s string, lookup map[string]string) (string, []string) {
	var missing []string
	seenMissing := map[string]bool{}
	result := expandVars(s, func(key string) string {
		if v, ok := lookup[key]; ok {
			return v
		}
		if !seenMissing[key] {
			missing = append(missing, key)
			seenMissing[key] = true
		}
		return fmt.Sprintf("${%s}", key)
	})
	return result, missing
}

// expandVars walks s and calls fn for every $KEY or ${KEY} reference.
func expandVars(s string, fn func(string) string) string {
	var b strings.Builder
	i := 0
	for i < len(s) {
		if s[i] != '$' || i+1 >= len(s) {
			b.WriteByte(s[i])
			i++
			continue
		}
		i++ // skip '$'
		if s[i] == '{' {
			i++ // skip '{'
			j := strings.IndexByte(s[i:], '}')
			if j < 0 {
				b.WriteString("${") // malformed, emit as-is
				continue
			}
			key := s[i : i+j]
			b.WriteString(fn(key))
			i += j + 1
		} else {
			j := i
			for j < len(s) && isVarChar(s[j]) {
				j++
			}
			if j == i {
				b.WriteByte('$')
				continue
			}
			b.WriteString(fn(s[i:j]))
			i = j
		}
	}
	return b.String()
}

func isVarChar(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') || c == '_'
}
