package envfile

import (
	"testing"
)

func makeEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestValidate_ValidEntries(t *testing.T) {
	entries := makeEntries("APP_NAME", "envoy", "PORT", "8080")
	result := Validate(entries)
	if !result.IsValid() {
		t.Fatalf("expected valid, got errors: %v", result.Errors)
	}
}

func TestValidate_EmptyKey(t *testing.T) {
	entries := makeEntries("", "value")
	result := Validate(entries)
	if result.IsValid() {
		t.Fatal("expected error for empty key")
	}
	if result.Errors[0].Message != "empty key is not allowed" {
		t.Errorf("unexpected message: %s", result.Errors[0].Message)
	}
}

func TestValidate_KeyStartsWithDigit(t *testing.T) {
	entries := makeEntries("1BAD_KEY", "value")
	result := Validate(entries)
	if result.IsValid() {
		t.Fatal("expected error for key starting with digit")
	}
}

func TestValidate_KeyWithInvalidChar(t *testing.T) {
	entries := makeEntries("BAD-KEY", "value")
	result := Validate(entries)
	if result.IsValid() {
		t.Fatal("expected error for key with hyphen")
	}
}

func TestValidate_DuplicateKeys(t *testing.T) {
	entries := makeEntries("FOO", "bar", "FOO", "baz")
	result := Validate(entries)
	if result.IsValid() {
		t.Fatal("expected error for duplicate key")
	}
	found := false
	for _, e := range result.Errors {
		if e.Key == "FOO" && e.Line == 2 {
			found = true
		}
	}
	if !found {
		t.Error("expected duplicate key error on line 2 for FOO")
	}
}

func TestValidate_SensitiveEmptyValue(t *testing.T) {
	entries := makeEntries("SECRET_KEY", "")
	result := Validate(entries)
	if result.IsValid() {
		t.Fatal("expected error for sensitive key with empty value")
	}
	if result.Errors[0].Message != "sensitive key has an empty value" {
		t.Errorf("unexpected message: %s", result.Errors[0].Message)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	entries := makeEntries("BAD-KEY", "val", "BAD-KEY", "val2")
	result := Validate(entries)
	if len(result.Errors) < 2 {
		t.Errorf("expected at least 2 errors, got %d", len(result.Errors))
	}
}
