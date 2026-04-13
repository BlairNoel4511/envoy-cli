package envfile

import (
	"encoding/json"
	"os"
)

// SaveGroupStore writes a GroupStore to a JSON file at the given path.
func SaveGroupStore(path string, gs *GroupStore) error {
	data, err := json.MarshalIndent(gs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadGroupStore reads a GroupStore from a JSON file.
// If the file does not exist an empty store is returned.
func LoadGroupStore(path string) (*GroupStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewGroupStore(), nil
		}
		return nil, err
	}
	gs := NewGroupStore()
	if err := json.Unmarshal(data, gs); err != nil {
		return nil, err
	}
	return gs, nil
}
