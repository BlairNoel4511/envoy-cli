package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// templateVarRe matches ${VAR_NAME} or $VAR_NAME style placeholders.
var templateVarRe = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// TemplateResult holds the outcome of a template expansion.
type TemplateResult struct {
	Expanded  string
	Missing   []string
	Resolved  []string
}

// ExpandTemplate replaces variable placeholders in a string using the
// provided entries as a lookup table. Missing variables are collected
// rather than causing an error, allowing callers to decide behaviour.
func ExpandTemplate(tmpl string, entries []Entry) TemplateResult {
	lookup := ToMap(entries)
	missing := []string{}
	resolved := []string{}
	seen := map[string]bool{}

	expanded := templateVarRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		key := extractKey(match)
		if val, ok := lookup[key]; ok {
			if !seen[key] {
				resolved = append(resolved, key)
				seen[key] = true
			}
			return val
		}
		if !seen[key] {
			missing = append(missing, key)
			seen[key] = true
		}
		return match
	})

	return TemplateResult{
		Expanded: expanded,
		Missing:  missing,
		Resolved: resolved,
	}
}

// ExpandEntries applies template expansion to every value in entries,
// using the full entry slice as the variable source (self-referential
// expansion is intentionally single-pass to avoid cycles).
func ExpandEntries(entries []Entry) ([]Entry, []string) {
	allMissing := []string{}
	result := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if !strings.Contains(e.Value, "$") {
			result = append(result, e)
			continue
		}
		r := ExpandTemplate(e.Value, entries)
		result = append(result, Entry{Key: e.Key, Value: r.Expanded})
		allMissing = append(allMissing, r.Missing...)
	}

	return result, allMissing
}

// FormatTemplateResult returns a human-readable summary of expansion.
func FormatTemplateResult(r TemplateResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Expanded: %s\n", r.Expanded))
	if len(r.Resolved) > 0 {
		sb.WriteString(fmt.Sprintf("Resolved: %s\n", strings.Join(r.Resolved, ", ")))
	}
	if len(r.Missing) > 0 {
		sb.WriteString(fmt.Sprintf("Missing:  %s\n", strings.Join(r.Missing, ", ")))
	}
	return sb.String()
}

func extractKey(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}
