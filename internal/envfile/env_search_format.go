package envfile

import (
	"fmt"
	"strings"
)

// FormatSearchResults formats search results for display.
func FormatSearchResults(results []SearchResult, colorize bool) string {
	if len(results) == 0 {
		return "  (no matches)"
	}

	var sb strings.Builder
	for _, r := range results {
		key := r.Entry.Key
		val := r.Entry.Value

		var tag string
		switch {
		case r.MatchedKey && r.MatchedValue:
			tag = "key+value"
		case r.MatchedKey:
			tag = "key"
		default:
			tag = "value"
		}

		line := fmt.Sprintf("  [%s] %s=%s", tag, key, val)
		if colorize {
			line = "\033[33m" + line + "\033[0m"
		}
		sb.WriteString(line + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// FormatSearchSummary returns a one-line summary of search results.
func FormatSearchSummary(s SearchSummary) string {
	return fmt.Sprintf("Search %q: %d/%d matched", s.Query, s.Matched, s.Total)
}
