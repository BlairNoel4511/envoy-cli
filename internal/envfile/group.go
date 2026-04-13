package envfile

import (
	"fmt"
	"sort"
)

// Group represents a named collection of env entry keys.
type Group struct {
	Name string   `json:"name"`
	Keys []string `json:"keys"`
}

// GroupStore holds multiple named groups.
type GroupStore struct {
	Groups map[string]*Group `json:"groups"`
}

// NewGroupStore initialises an empty GroupStore.
func NewGroupStore() *GroupStore {
	return &GroupStore{Groups: make(map[string]*Group)}
}

// Add creates or updates a group with the given name and keys.
func (gs *GroupStore) Add(name string, keys []string) {
	if _, ok := gs.Groups[name]; !ok {
		gs.Groups[name] = &Group{Name: name}
	}
	existing := gs.Groups[name]
	seen := make(map[string]bool)
	for _, k := range existing.Keys {
		seen[k] = true
	}
	for _, k := range keys {
		if !seen[k] {
			existing.Keys = append(existing.Keys, k)
			seen[k] = true
		}
	}
}

// Remove deletes a group by name. Returns an error if not found.
func (gs *GroupStore) Remove(name string) error {
	if _, ok := gs.Groups[name]; !ok {
		return fmt.Errorf("group %q not found", name)
	}
	delete(gs.Groups, name)
	return nil
}

// Get returns the group with the given name, or nil if absent.
func (gs *GroupStore) Get(name string) *Group {
	return gs.Groups[name]
}

// List returns all group names in sorted order.
func (gs *GroupStore) List() []string {
	names := make([]string, 0, len(gs.Groups))
	for n := range gs.Groups {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// FilterByGroup returns only the entries whose keys belong to the named group.
func FilterByGroup(entries []Entry, gs *GroupStore, name string) ([]Entry, error) {
	g := gs.Get(name)
	if g == nil {
		return nil, fmt.Errorf("group %q not found", name)
	}
	allowed := make(map[string]bool, len(g.Keys))
	for _, k := range g.Keys {
		allowed[k] = true
	}
	var result []Entry
	for _, e := range entries {
		if allowed[e.Key] {
			result = append(result, e)
		}
	}
	return result, nil
}
