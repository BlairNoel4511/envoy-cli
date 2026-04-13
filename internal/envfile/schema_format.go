package envfile

import (
	"fmt"
	"strings"
)

// FormatSchemaViolations returns a human-readable summary of schema violations.
func FormatSchemaViolations(violations []SchemaViolation, colorize bool) string {
	if len(violations) == 0 {
		if colorize {
			return "\033[32m✔ schema validation passed\033[0m"
		}
		return "✔ schema validation passed"
	}

	var sb strings.Builder
	for _, v := range violations {
		if colorize {
			sb.WriteString(fmt.Sprintf("\033[31m✘ [%s]\033[0m %s\n", v.Key, v.Message))
		} else {
			sb.WriteString(fmt.Sprintf("✘ [%s] %s\n", v.Key, v.Message))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// FormatSchemaSummaryLine returns a one-line summary of the validation result.
func FormatSchemaSummaryLine(violations []SchemaViolation) string {
	if len(violations) == 0 {
		return "schema: 0 violations"
	}
	return fmt.Sprintf("schema: %d violation(s) found", len(violations))
}

// ViolationKeys returns just the keys that have violations.
func ViolationKeys(violations []SchemaViolation) []string {
	seen := make(map[string]struct{})
	var keys []string
	for _, v := range violations {
		if _, ok := seen[v.Key]; !ok {
			seen[v.Key] = struct{}{}
			keys = append(keys, v.Key)
		}
	}
	return keys
}
