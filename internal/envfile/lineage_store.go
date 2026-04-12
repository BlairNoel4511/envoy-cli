package envfile

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// SaveLineage writes a Lineage to a JSON file at the given path.
func SaveLineage(path string, l *Lineage) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// LoadLineage reads a Lineage from a JSON file at the given path.
// If the file does not exist, an empty Lineage is returned.
func LoadLineage(path string) (*Lineage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return NewLineage(), nil
		}
		return nil, err
	}
	var l Lineage
	if err := json.Unmarshal(data, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// AppendLineageEvent loads a lineage from path, appends a single event, and saves it back.
func AppendLineageEvent(path, action, key, oldValue, newValue, source string) error {
	l, err := LoadLineage(path)
	if err != nil {
		return err
	}
	l.Record(action, key, oldValue, newValue, source)
	return SaveLineage(path, l)
}
