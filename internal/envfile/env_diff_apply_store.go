package envfile

import (
	"encoding/json"
	"os"
)

// SaveDiffApplyResults persists a slice of DiffApplyResult to a JSON file.
func SaveDiffApplyResults(path string, results []DiffApplyResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadDiffApplyResults reads a previously saved DiffApplyResult slice from disk.
// If the file does not exist, an empty slice is returned without error.
func LoadDiffApplyResults(path string) ([]DiffApplyResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []DiffApplyResult{}, nil
		}
		return nil, err
	}
	var results []DiffApplyResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results, nil
}
