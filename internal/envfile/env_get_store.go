package envfile

import (
	"encoding/json"
	"os"
)

// SaveGetResults persists get results to a JSON file.
func SaveGetResults(path string, results []GetResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadGetResults reads previously saved get results from a JSON file.
func LoadGetResults(path string) ([]GetResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []GetResult{}, nil
		}
		return nil, err
	}
	var results []GetResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results, nil
}
