package envfile

import (
	"strings"
	"testing"
)

func makeSchemaEntries() []Entry {
	return []Entry{
		{Key: "PORT", Value: "8080"},
		{Key: "DEBUG", Value: "true"},
		{Key: "API_URL", Value: "https://example.com"},
		{Key: "APP_NAME", Value: "envoy"},
	}
}

func TestValidateSchema_NoViolations(t *testing.T) {
	entries := makeSchemaEntries()
	schema := Schema{
		Fields: []SchemaField{
			{Key: "PORT", Type: FieldTypeInt, Required: true},
			{Key: "DEBUG", Type: FieldTypeBool, Required: true},
			{Key: "API_URL", Type: FieldTypeURL, Required: true},
		},
	}
	violations := ValidateSchema(entries, schema)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d: %+v", len(violations), violations)
	}
}

func TestValidateSchema_MissingRequiredKey(t *testing.T) {
	entries := []Entry{{Key: "APP_NAME", Value: "envoy"}}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "PORT", Type: FieldTypeInt, Required: true},
		},
	}
	violations := ValidateSchema(entries, schema)
	if len(violations) != 1 || violations[0].Key != "PORT" {
		t.Errorf("expected violation for PORT, got %+v", violations)
	}
}

func TestValidateSchema_InvalidInt(t *testing.T) {
	entries := []Entry{{Key: "PORT", Value: "not-a-number"}}
	schema := Schema{
		Fields: []SchemaField{{Key: "PORT", Type: FieldTypeInt}},
	}
	violations := ValidateSchema(entries, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "expected int") {
		t.Errorf("unexpected message: %s", violations[0].Message)
	}
}

func TestValidateSchema_InvalidBool(t *testing.T) {
	entries := []Entry{{Key: "DEBUG", Value: "maybe"}}
	schema := Schema{
		Fields: []SchemaField{{Key: "DEBUG", Type: FieldTypeBool}},
	}
	violations := ValidateSchema(entries, schema)
	if len(violations) != 1 || !strings.Contains(violations[0].Message, "expected bool") {
		t.Errorf("expected bool violation, got %+v", violations)
	}
}

func TestValidateSchema_InvalidURL(t *testing.T) {
	entries := []Entry{{Key: "API_URL", Value: "not-a-url"}}
	schema := Schema{
		Fields: []SchemaField{{Key: "API_URL", Type: FieldTypeURL}},
	}
	violations := ValidateSchema(entries, schema)
	if len(violations) != 1 || !strings.Contains(violations[0].Message, "expected URL") {
		t.Errorf("expected URL violation, got %+v", violations)
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	entries := []Entry{{Key: "ENV", Value: "staging"}}
	schema := Schema{
		Fields: []SchemaField{{Key: "ENV", Type: FieldTypeString, Pattern: "^(production|development)$"}},
	}
	violations := ValidateSchema(entries, schema)
	if len(violations) != 1 || !strings.Contains(violations[0].Message, "does not match pattern") {
		t.Errorf("expected pattern violation, got %+v", violations)
	}
}

func TestFormatSchemaViolations_NoViolations(t *testing.T) {
	out := FormatSchemaViolations(nil, false)
	if !strings.Contains(out, "passed") {
		t.Errorf("expected pass message, got: %s", out)
	}
}

func TestFormatSchemaViolations_WithViolations(t *testing.T) {
	v := []SchemaViolation{{Key: "PORT", Message: "required key is missing"}}
	out := FormatSchemaViolations(v, false)
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got: %s", out)
	}
}

func TestFormatSchemaSummaryLine(t *testing.T) {
	if s := FormatSchemaSummaryLine(nil); !strings.Contains(s, "0 violations") {
		t.Errorf("unexpected: %s", s)
	}
	v := []SchemaViolation{{Key: "X", Message: "bad"}}
	if s := FormatSchemaSummaryLine(v); !strings.Contains(s, "1 violation") {
		t.Errorf("unexpected: %s", s)
	}
}
