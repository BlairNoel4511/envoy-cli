package envfile

import (
	"encoding/json"
	"os"
)

// SavePinStore writes the PinStore to disk as JSON.
func SavePinStore(path string, store *PinStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadPinStore reads a PinStore from disk.
// If the file does not exist, an empty store is returned.
func LoadPinStore(path string) (*PinStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewPinStore(), nil
		}
		return nil, err
	}
	var store PinStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	if store.Pins == nil {
		store.Pins = []PinEntry{}
	}
	return &store, nil
}
