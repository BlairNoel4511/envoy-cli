package envfile

import (
	"encoding/json"
	"os"
)

// ProtectStore tracks which keys are protected.
type ProtectStore struct {
	keys map[string]struct{}
}

// NewProtectStore returns an empty ProtectStore.
func NewProtectStore() *ProtectStore {
	return &ProtectStore{keys: make(map[string]struct{})}
}

// Add marks a key as protected.
func (p *ProtectStore) Add(key string) {
	p.keys[key] = struct{}{}
}

// Remove unprotects a key.
func (p *ProtectStore) Remove(key string) {
	delete(p.keys, key)
}

// IsProtected returns true if the key is currently protected.
func (p *ProtectStore) IsProtected(key string) bool {
	_, ok := p.keys[key]
	return ok
}

// All returns all protected keys as a slice.
func (p *ProtectStore) All() []string {
	out := make([]string, 0, len(p.keys))
	for k := range p.keys {
		out = append(out, k)
	}
	return out
}

// SaveProtectStore persists the store to a JSON file.
func SaveProtectStore(path string, store *ProtectStore) error {
	data, err := json.MarshalIndent(store.All(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadProtectStore loads a ProtectStore from a JSON file.
// Returns an empty store if the file does not exist.
func LoadProtectStore(path string) (*ProtectStore, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NewProtectStore(), nil
	}
	if err != nil {
		return nil, err
	}
	var keys []string
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, err
	}
	s := NewProtectStore()
	for _, k := range keys {
		s.Add(k)
	}
	return s, nil
}
