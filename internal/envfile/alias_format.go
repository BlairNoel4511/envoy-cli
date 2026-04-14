package envfile

import (
	"fmt"
	"strings"
)

// FormatAliasList returns a human-readable list of all aliases.
func FormatAliasList(aliases []Alias) string {
	if len(aliases) == 0 {
		return "(no aliases defined)\n"
	}
	var sb strings.Builder
	for _, a := range aliases {
		line := fmt.Sprintf("  %-24s -> %s", a.Alias, a.Canonical)
		if a.Comment != "" {
			line += fmt.Sprintf("  # %s", a.Comment)
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}

// FormatAliasDetail returns a detailed view of a single alias entry.
func FormatAliasDetail(a Alias) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Alias:     %s\n", a.Alias))
	sb.WriteString(fmt.Sprintf("Canonical: %s\n", a.Canonical))
	if a.Comment != "" {
		sb.WriteString(fmt.Sprintf("Comment:   %s\n", a.Comment))
	}
	return sb.String()
}

// FormatResolvedEntries shows which alias keys were resolved to canonical keys.
func FormatResolvedEntries(original, resolved []Entry) string {
	if len(original) == 0 {
		return "(no entries)\n"
	}
	var sb strings.Builder
	for i, orig := range original {
		if i >= len(resolved) {
			break
		}
		res := resolved[i]
		if orig.Key != res.Key {
			sb.WriteString(fmt.Sprintf("  %s -> %s\n", orig.Key, res.Key))
		} else {
			sb.WriteString(fmt.Sprintf("  %s (unchanged)\n", orig.Key))
		}
	}
	return sb.String()
}
