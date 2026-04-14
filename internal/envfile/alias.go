package envfile

import "fmt"

// Alias represents a mapping from an alias name to a canonical key.
type Alias struct {
	Alias    string `json:"alias"`
	Canonical string `json:"canonical"`
	Comment  string `json:"comment,omitempty"`
}

// AliasStore holds alias mappings for env keys.
type AliasStore struct {
	aliases map[string]Alias // alias -> Alias
}

// NewAliasStore creates an empty AliasStore.
func NewAliasStore() *AliasStore {
	return &AliasStore{aliases: make(map[string]Alias)}
}

// Add registers an alias pointing to a canonical key.
func (s *AliasStore) Add(alias, canonical, comment string) error {
	if alias == "" {
		return fmt.Errorf("alias name must not be empty")
	}
	if canonical == "" {
		return fmt.Errorf("canonical key must not be empty")
	}
	if alias == canonical {
		return fmt.Errorf("alias %q cannot point to itself", alias)
	}
	s.aliases[alias] = Alias{Alias: alias, Canonical: canonical, Comment: comment}
	return nil
}

// Get returns the Alias entry for the given alias name.
func (s *AliasStore) Get(alias string) (Alias, bool) {
	a, ok := s.aliases[alias]
	return a, ok
}

// Remove deletes an alias.
func (s *AliasStore) Remove(alias string) bool {
	_, ok := s.aliases[alias]
	if ok {
		delete(s.aliases, alias)
	}
	return ok
}

// List returns all aliases sorted by alias name.
func (s *AliasStore) List() []Alias {
	out := make([]Alias, 0, len(s.aliases))
	for _, a := range s.aliases {
		out = append(out, a)
	}
	sortAliases(out)
	return out
}

// Resolve returns the canonical key for an alias, or the original key if no alias exists.
func (s *AliasStore) Resolve(key string) string {
	if a, ok := s.aliases[key]; ok {
		return a.Canonical
	}
	return key
}

// ResolveEntries returns a copy of entries with alias keys replaced by their canonical keys.
func (s *AliasStore) ResolveEntries(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		e.Key = s.Resolve(e.Key)
		out[i] = e
	}
	return out
}

func sortAliases(a []Alias) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j].Alias < a[j-1].Alias; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}
