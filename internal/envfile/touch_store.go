package envfile

import (
	"encoding/json"
	"os"
	"time"
)

// TouchStore holds per-key timestamps recording when a key was last touched.
type TouchStore struct {
	entries map[string]time.Time
}

// NewTouchStore returns an empty TouchStore.
func NewTouchStore() *TouchStore {
	return &TouchStore{entries: make(map[string]time.Time)}
}

// Set records the timestamp for key.
func (ts *TouchStore) Set(key string, t time.Time) {
	ts.entries[key] = t
}

// Get returns the timestamp for key and whether it existed.
func (ts *TouchStore) Get(key string) (time.Time, bool) {
	t, ok := ts.entries[key]
	return t, ok
}

// Remove deletes the timestamp for key.
func (ts *TouchStore) Remove(key string) {
	delete(ts.entries, key)
}

// All returns a copy of all key→timestamp pairs.
func (ts *TouchStore) All() map[string]time.Time {
	out := make(map[string]time.Time, len(ts.entries))
	for k, v := range ts.entries {
		out[k] = v
	}
	return out
}

// SaveTouchStore persists the store to path as JSON with restricted permissions.
func SaveTouchStore(path string, ts *TouchStore) error {
	data, err := json.MarshalIndent(ts.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadTouchStore reads a TouchStore from path.
// If the file does not exist an empty store is returned without error.
func LoadTouchStore(path string) (*TouchStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewTouchStore(), nil
		}
		return nil, err
	}
	m := make(map[string]time.Time)
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &TouchStore{entries: m}, nil
}
