package envfile

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidationError represents a single validation issue found in an env file.
type ValidationError struct {
	Line    int
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("line %d: %s", e.Line, e.Message)
	}
	return e.Message
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

func (r *ValidationResult) Add(err ValidationError) {
	r.Errors = append(r.Errors, err)
}

// Validate checks an env map for common issues such as empty keys,
// keys with invalid characters, or suspiciously empty values for sensitive vars.
func Validate(entries []Entry) ValidationResult {
	result := ValidationResult{}
	seen := make(map[string]int)

	for i, entry := range entries {
		lineNum := i + 1

		if entry.Key == "" {
			result.Add(ValidationError{Line: lineNum, Message: "empty key is not allowed"})
			continue
		}

		if err := validateKey(entry.Key); err != nil {
			result.Add(ValidationError{Line: lineNum, Key: entry.Key, Message: err.Error()})
		}

		if prev, ok := seen[entry.Key]; ok {
			result.Add(ValidationError{
				Line:    lineNum,
				Key:     entry.Key,
				Message: fmt.Sprintf("duplicate key (first seen on line %d)", prev),
			})
		} else {
			seen[entry.Key] = lineNum
		}

		if IsSensitive(entry.Key) && strings.TrimSpace(entry.Value) == "" {
			result.Add(ValidationError{
				Line:    lineNum,
				Key:     entry.Key,
				Message: "sensitive key has an empty value",
			})
		}
	}

	return result
}

func validateKey(key string) error {
	for i, ch := range key {
		if i == 0 && unicode.IsDigit(ch) {
			return fmt.Errorf("key %q must not start with a digit", key)
		}
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return fmt.Errorf("key %q contains invalid character %q", key, ch)
		}
	}
	return nil
}
