package envfile

import (
	"encoding/json"
	"os"
	"time"
)

// FreezeStore tracks which keys are frozen and when they were frozen.
type FreezeStore struct {
	Entries map[string]time.Time `json:"entries"`
}

// NewFreezeStore returns an empty FreezeStore.
func NewFreezeStore() *FreezeStore {
	return &FreezeStore{Entries: make(map[string]time.Time)}
}

// Freeze marks key as frozen at the given time.
func (f *FreezeStore) Freeze(key string, at time.Time) {
	f.Entries[key] = at
}

// Unfreeze removes the frozen status for key.
func (f *FreezeStore) Unfreeze(key string) {
	delete(f.Entries, key)
}

// IsFrozen reports whether key is currently frozen.
func (f *FreezeStore) IsFrozen(key string) bool {
	_, ok := f.Entries[key]
	return ok
}

// FrozenAt returns the time at which key was frozen and whether it exists.
func (f *FreezeStore) FrozenAt(key string) (time.Time, bool) {
	t, ok := f.Entries[key]
	return t, ok
}

// SaveFreezeStore serialises the store to path with restricted permissions.
func SaveFreezeStore(path string, store *FreezeStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadFreezeStore deserialises a FreezeStore from path.
// If the file does not exist an empty store is returned.
func LoadFreezeStore(path string) (*FreezeStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewFreezeStore(), nil
		}
		return nil, err
	}
	var store FreezeStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	if store.Entries == nil {
		store.Entries = make(map[string]time.Time)
	}
	return &store, nil
}
