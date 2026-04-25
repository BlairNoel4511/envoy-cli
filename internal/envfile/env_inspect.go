package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// InspectResult holds metadata about a single env entry.
type InspectResult struct {
	Key       string
	Value     string
	Redacted  bool
	Sensitive bool
	HasComment bool
	Comment   string
	Length    int
	Found     bool
}

// InspectOptions controls how inspection behaves.
type InspectOptions struct {
	Redact bool
}

// Inspect returns detailed metadata for the requested keys.
// If keys is empty, all entries are inspected.
func Inspect(entries []Entry, keys []string, opts InspectOptions) []InspectResult {
	byKey := make(map[string]Entry, len(entries))
	for _, e := range entries {
		byKey[e.Key] = e
	}

	targets := keys
	if len(targets) == 0 {
		for _, e := range entries {
			targets = append(targets, e.Key)
		}
		sort.Strings(targets)
	}

	results := make([]InspectResult, 0, len(targets))
	for _, k := range targets {
		e, ok := byKey[k]
		if !ok {
			results = append(results, InspectResult{Key: k, Found: false})
			continue
		}
		sensitive := IsSensitive(k)
		value := e.Value
		redacted := false
		if opts.Redact && sensitive {
			value = "***"
			redacted = true
		}
		results = append(results, InspectResult{
			Key:        k,
			Value:      value,
			Redacted:   redacted,
			Sensitive:  sensitive,
			HasComment: strings.TrimSpace(e.Comment) != "",
			Comment:    strings.TrimSpace(e.Comment),
			Length:     len(e.Value),
			Found:      true,
		})
	}
	return results
}

// InspectSummary returns a human-readable summary line.
func InspectSummary(results []InspectResult) string {
	total, found, sensitive := 0, 0, 0
	for _, r := range results {
		total++
		if r.Found {
			found++
		}
		if r.Sensitive {
			sensitive++
		}
	}
	return fmt.Sprintf("%d inspected, %d found, %d sensitive", total, found, sensitive)
}
