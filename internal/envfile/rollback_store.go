package envfile

import (
	"encoding/json"
	"os"
)

// SaveRollbackPlan persists a RollbackPlan to a JSON file at the given path.
func SaveRollbackPlan(path string, plan RollbackPlan) error {
	data, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadRollbackPlan reads a RollbackPlan from a JSON file at the given path.
// Returns an empty plan and no error if the file does not exist.
func LoadRollbackPlan(path string) (RollbackPlan, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return RollbackPlan{}, nil
		}
		return RollbackPlan{}, err
	}

	var plan RollbackPlan
	if err := json.Unmarshal(data, &plan); err != nil {
		return RollbackPlan{}, err
	}
	return plan, nil
}
