package envfile

import (
	"encoding/json"
	"os"
)

type tagStoreFile struct {
	Tags map[string][]string `json:"tags"`
}

// SaveTagStore writes a TagStore to disk as JSON.
func SaveTagStore(path string, ts *TagStore) error {
	data, err := json.MarshalIndent(&tagStoreFile{Tags: ts.tags}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadTagStore reads a TagStore from a JSON file.
// If the file does not exist, an empty store is returned.
func LoadTagStore(path string) (*TagStore, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NewTagStore(), nil
	}
	if err != nil {
		return nil, err
	}
	var f tagStoreFile
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	ts := NewTagStore()
	if f.Tags != nil {
		ts.tags = f.Tags
	}
	return ts, nil
}
