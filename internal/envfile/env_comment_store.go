package envfile

import (
	"encoding/json"
	"os"
)

// SaveCommentResults persists comment results to a JSON file.
func SaveCommentResults(path string, results []CommentResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadCommentResults loads comment results from a JSON file.
// Returns an empty slice if the file does not exist.
func LoadCommentResults(path string) ([]CommentResult, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []CommentResult{}, nil
	}
	if err != nil {
		return nil, err
	}
	var results []CommentResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results, nil
}
