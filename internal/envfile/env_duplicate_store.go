package envfile

import (
	"encoding/json"
	"os"
)

// SaveDuplicateResults persists a slice of DuplicateResult to a JSON file.
func SaveDuplicateResults(path string, results []DuplicateResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadDuplicateResults reads a slice of DuplicateResult from a JSON file.
// If the file does not exist, an empty slice is returned without error.
func LoadDuplicateResults(path string) ([]DuplicateResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []DuplicateResult{}, nil
		}
		return nil, err
	}
	var results []DuplicateResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results, nil
}
