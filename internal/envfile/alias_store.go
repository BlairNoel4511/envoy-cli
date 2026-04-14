package envfile

import (
	"encoding/json"
	"os"
)

// SaveAliasStore persists the AliasStore to a JSON file.
func SaveAliasStore(path string, s *AliasStore) error {
	data, err := json.MarshalIndent(s.aliases, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadAliasStore reads an AliasStore from a JSON file.
// Returns an empty store if the file does not exist.
func LoadAliasStore(path string) (*AliasStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewAliasStore(), nil
		}
		return nil, err
	}
	var raw map[string]Alias
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	s := NewAliasStore()
	s.aliases = raw
	return s, nil
}
