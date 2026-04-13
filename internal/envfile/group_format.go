package envfile

import (
	"fmt"
	"strings"
)

// FormatGroupList returns a human-readable summary of all groups in the store.
func FormatGroupList(gs *GroupStore) string {
	names := gs.List()
	if len(names) == 0 {
		return "No groups defined."
	}
	var sb strings.Builder
	for _, name := range names {
		g := gs.Get(name)
		sb.WriteString(fmt.Sprintf("[%s] (%d keys)\n", g.Name, len(g.Keys)))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// FormatGroupDetail returns the keys belonging to a single group.
func FormatGroupDetail(g *Group) string {
	if g == nil {
		return "Group not found."
	}
	if len(g.Keys) == 0 {
		return fmt.Sprintf("Group %q is empty.", g.Name)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Group: %s\n", g.Name))
	for _, k := range g.Keys {
		sb.WriteString(fmt.Sprintf("  - %s\n", k))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// FormatFilteredByGroup returns a display of entries filtered by group.
func FormatFilteredByGroup(entries []Entry, groupName string) string {
	if len(entries) == 0 {
		return fmt.Sprintf("No entries matched group %q.", groupName)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Entries in group %q:\n", groupName))
	for _, e := range entries {
		val := e.Value
		if IsSensitive(e.Key) {
			val = "***"
		}
		sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, val))
	}
	return strings.TrimRight(sb.String(), "\n")
}
