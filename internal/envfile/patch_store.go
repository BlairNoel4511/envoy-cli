package envfile

import (
	"encoding/json"
	"os"
)

// SavePatchInstructions writes a list of patch instructions to a JSON file.
func SavePatchInstructions(path string, instructions []PatchInstruction) error {
	data, err := json.MarshalIndent(instructions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadPatchInstructions reads patch instructions from a JSON file.
// Returns an empty slice if the file does not exist.
func LoadPatchInstructions(path string) ([]PatchInstruction, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []PatchInstruction{}, nil
		}
		return nil, err
	}
	var instructions []PatchInstruction
	if err := json.Unmarshal(data, &instructions); err != nil {
		return nil, err
	}
	return instructions, nil
}
