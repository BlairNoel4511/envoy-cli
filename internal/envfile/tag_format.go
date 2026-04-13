package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatTagList returns a human-readable listing of all tags.
func FormatTagList(ts *TagStore) string {
	all := ts.All()
	if len(all) == 0 {
		return "(no tags)"
	}
	// group by key
	grouped := make(map[string][]string)
	for _, t := range all {
		grouped[t.Key] = append(grouped[t.Key], t.Label)
	}
	keys := make([]string, 0, len(grouped))
	for k := range grouped {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		labels := grouped[k]
		sort.Strings(labels)
		sb.WriteString(fmt.Sprintf("  %-24s %s\n", k, strings.Join(labels, ", ")))
	}
	return sb.String()
}

// FormatTagsForKey returns a single-line summary of tags for a key.
func FormatTagsForKey(key string, labels []string) string {
	if len(labels) == 0 {
		return fmt.Sprintf("%s: (no tags)", key)
	}
	sorted := make([]string, len(labels))
	copy(sorted, labels)
	sort.Strings(sorted)
	return fmt.Sprintf("%s: %s", key, strings.Join(sorted, ", "))
}

// FormatKeysWithTag returns a summary of keys sharing a label.
func FormatKeysWithTag(label string, keys []string) string {
	if len(keys) == 0 {
		return fmt.Sprintf("No keys tagged %q", label)
	}
	sorted := make([]string, len(keys))
	copy(sorted, keys)
	sort.Strings(sorted)
	return fmt.Sprintf("Keys tagged %q:\n  %s", label, strings.Join(sorted, "\n  "))
}
