package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaFieldType represents the expected type for a schema field.
type SchemaFieldType string

const (
	FieldTypeString SchemaFieldType = "string"
	FieldTypeInt    SchemaFieldType = "int"
	FieldTypeBool   SchemaFieldType = "bool"
	FieldTypeURL    SchemaFieldType = "url"
)

// SchemaField defines constraints for a single env key.
type SchemaField struct {
	Key      string
	Type     SchemaFieldType
	Required bool
	Pattern  string // optional regex pattern
}

// Schema holds a collection of field definitions.
type Schema struct {
	Fields []SchemaField
}

// SchemaViolation describes a single validation failure.
type SchemaViolation struct {
	Key     string
	Message string
}

var (
	reInt  = regexp.MustCompile(`^-?\d+$`)
	reBool = regexp.MustCompile(`^(?i)(true|false|1|0|yes|no)$`)
	reURL  = regexp.MustCompile(`^https?://`)
)

// ValidateSchema checks a set of entries against the schema and returns violations.
func ValidateSchema(entries []Entry, schema Schema) []SchemaViolation {
	entryMap := ToMap(entries)
	var violations []SchemaViolation

	for _, field := range schema.Fields {
		val, exists := entryMap[field.Key]

		if field.Required && !exists {
			violations = append(violations, SchemaViolation{
				Key:     field.Key,
				Message: "required key is missing",
			})
			continue
		}

		if !exists {
			continue
		}

		switch field.Type {
		case FieldTypeInt:
			if !reInt.MatchString(val) {
				violations = append(violations, SchemaViolation{Key: field.Key, Message: fmt.Sprintf("expected int, got %q", val)})
			}
		case FieldTypeBool:
			if !reBool.MatchString(val) {
				violations = append(violations, SchemaViolation{Key: field.Key, Message: fmt.Sprintf("expected bool, got %q", val)})
			}
		case FieldTypeURL:
			if !reURL.MatchString(strings.TrimSpace(val)) {
				violations = append(violations, SchemaViolation{Key: field.Key, Message: fmt.Sprintf("expected URL (http/https), got %q", val)})
			}
		}

		if field.Pattern != "" {
			re, err := regexp.Compile(field.Pattern)
			if err == nil && !re.MatchString(val) {
				violations = append(violations, SchemaViolation{Key: field.Key, Message: fmt.Sprintf("value %q does not match pattern %q", val, field.Pattern)})
			}
		}
	}

	return violations
}
