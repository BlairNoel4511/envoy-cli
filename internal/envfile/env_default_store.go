package envfile

import (
	"encoding/json"
	"os"
)

// SaveDefaultResults persists a slice of DefaultResult to a JSON file.
func SaveDefaultResults(path string, results []DefaultResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadDefaultResults reads a previously saved DefaultResult slice from disk.
// Returns an empty slice when the file does not exist.
func LoadDefaultResults(path string) ([]DefaultResult, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []DefaultResult{}, nil
	}
	if err != nil {
		return nil, err
	}
	var results []DefaultResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results, nil
}
