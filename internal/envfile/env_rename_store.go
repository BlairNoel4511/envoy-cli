package envfile

import (
	"encoding/json"
	"os"
)

// RenameMapping is a serialisable map of old→new key names.
type RenameMapping struct {
	Mapping map[string]string `json:"mapping"`
}

// SaveRenameMapping persists a RenameMapping to a JSON file.
func SaveRenameMapping(path string, rm RenameMapping) error {
	data, err := json.MarshalIndent(rm, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadRenameMapping loads a RenameMapping from a JSON file.
// If the file does not exist an empty mapping is returned without error.
func LoadRenameMapping(path string) (RenameMapping, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return RenameMapping{Mapping: map[string]string{}}, nil
	}
	if err != nil {
		return RenameMapping{}, err
	}
	var rm RenameMapping
	if err := json.Unmarshal(data, &rm); err != nil {
		return RenameMapping{}, err
	}
	if rm.Mapping == nil {
		rm.Mapping = map[string]string{}
	}
	return rm, nil
}
