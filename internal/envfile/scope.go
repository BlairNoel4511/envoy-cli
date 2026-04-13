package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// Scope represents a named environment scope (e.g. "dev", "staging", "prod").
type Scope struct {
	Name    string            `json:"name"`
	Entries []Entry           `json:"entries"`
	Meta    map[string]string `json:"meta,omitempty"`
}

// ScopeStore holds multiple named scopes.
type ScopeStore struct {
	Scopes map[string]*Scope `json:"scopes"`
}

// NewScopeStore creates an empty ScopeStore.
func NewScopeStore() *ScopeStore {
	return &ScopeStore{Scopes: make(map[string]*Scope)}
}

// Set stores entries under a named scope, replacing any existing scope with that name.
func (s *ScopeStore) Set(name string, entries []Entry) {
	name = strings.ToLower(strings.TrimSpace(name))
	s.Scopes[name] = &Scope{
		Name:    name,
		Entries: entries,
		Meta:    make(map[string]string),
	}
}

// Get retrieves the entries for a named scope. Returns nil, false if not found.
func (s *ScopeStore) Get(name string) ([]Entry, bool) {
	name = strings.ToLower(strings.TrimSpace(name))
	sc, ok := s.Scopes[name]
	if !ok {
		return nil, false
	}
	return sc.Entries, true
}

// Remove deletes a scope by name.
func (s *ScopeStore) Remove(name string) {
	delete(s.Scopes, strings.ToLower(strings.TrimSpace(name)))
}

// List returns sorted scope names.
func (s *ScopeStore) List() []string {
	names := make([]string, 0, len(s.Scopes))
	for k := range s.Scopes {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// SetMeta attaches a metadata key/value to a scope.
func (s *ScopeStore) SetMeta(scope, key, value string) error {
	sc, ok := s.Scopes[strings.ToLower(scope)]
	if !ok {
		return fmt.Errorf("scope %q not found", scope)
	}
	sc.Meta[key] = value
	return nil
}

// FormatScopeList returns a human-readable summary of all scopes.
func FormatScopeList(store *ScopeStore) string {
	if store == nil || len(store.Scopes) == 0 {
		return "No scopes defined.\n"
	}
	var sb strings.Builder
	for _, name := range store.List() {
		sc := store.Scopes[name]
		sb.WriteString(fmt.Sprintf("  %-20s %d key(s)\n", sc.Name, len(sc.Entries)))
	}
	return sb.String()
}
