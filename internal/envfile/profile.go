package envfile

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Profile represents a named environment configuration (e.g. "dev", "staging", "prod").
type Profile struct {
	Name    string   `json:"name"`
	Entries []Entry  `json:"entries"`
	Tags    []string `json:"tags,omitempty"`
}

// ProfileStore holds multiple named profiles persisted to a single JSON file.
type ProfileStore struct {
	Profiles map[string]*Profile `json:"profiles"`
}

// NewProfileStore returns an empty ProfileStore.
func NewProfileStore() *ProfileStore {
	return &ProfileStore{Profiles: make(map[string]*Profile)}
}

// Set adds or replaces a profile in the store.
func (ps *ProfileStore) Set(p *Profile) {
	ps.Profiles[p.Name] = p
}

// Get returns a profile by name, or nil if not found.
func (ps *ProfileStore) Get(name string) (*Profile, bool) {
	p, ok := ps.Profiles[name]
	return p, ok
}

// Remove deletes a profile by name.
func (ps *ProfileStore) Remove(name string) {
	delete(ps.Profiles, name)
}

// List returns all profile names in the store.
func (ps *ProfileStore) List() []string {
	names := make([]string, 0, len(ps.Profiles))
	for name := range ps.Profiles {
		names = append(names, name)
	}
	return names
}

// SaveProfileStore writes the ProfileStore to the given file path as JSON.
func SaveProfileStore(path string, ps *ProfileStore) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(ps, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadProfileStore reads a ProfileStore from the given file path.
// Returns an empty store if the file does not exist.
func LoadProfileStore(path string) (*ProfileStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return NewProfileStore(), nil
		}
		return nil, err
	}
	var ps ProfileStore
	if err := json.Unmarshal(data, &ps); err != nil {
		return nil, err
	}
	if ps.Profiles == nil {
		ps.Profiles = make(map[string]*Profile)
	}
	return &ps, nil
}
