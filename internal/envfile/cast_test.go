package envfile

import (
	"testing"
)

func makeCastEntries() []Entry {
	return []Entry{
		{Key: "PORT", Value: "8080"},
		{Key: "RATE", Value: "3.14"},
		{Key: "DEBUG", Value: "1"},
		{Key: "NAME", Value: "envoy"},
		{Key: "SECRET", Value: "abc123"},
	}
}

func TestCast_IntConversion(t *testing.T) {
	entries := makeCastEntries()
	opts := CastOptions{Types: map[string]string{"PORT": "int"}}
	out, results, err := Cast(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "PORT" {
		t.Fatalf("expected 1 result for PORT, got %d", len(results))
	}
	if results[0].Casted != "8080" {
		t.Errorf("expected casted=8080, got %q", results[0].Casted)
	}
	if out[0].Value != "8080" {
		t.Errorf("expected out value=8080, got %q", out[0].Value)
	}
}

func TestCast_FloatConversion(t *testing.T) {
	entries := makeCastEntries()
	opts := CastOptions{Types: map[string]string{"RATE": "float"}}
	out, results, err := Cast(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Casted != "3.14" {
		t.Errorf("expected 3.14, got %q", results[0].Casted)
	}
	_ = out
}

func TestCast_BoolConversion(t *testing.T) {
	entries := makeCastEntries()
	opts := CastOptions{Types: map[string]string{"DEBUG": "bool"}}
	out, results, err := Cast(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Casted != "true" {
		t.Errorf("expected true, got %q", results[0].Casted)
	}
	if out[2].Value != "true" {
		t.Errorf("expected entry updated to true, got %q", out[2].Value)
	}
}

func TestCast_StringNoChange(t *testing.T) {
	entries := makeCastEntries()
	opts := CastOptions{Types: map[string]string{"NAME": "string"}}
	_, results, err := Cast(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Changed {
		t.Error("expected no change for string cast")
	}
}

func TestCast_InvalidIntReturnsError(t *testing.T) {
	entries := makeCastEntries()
	opts := CastOptions{Types: map[string]string{"NAME": "int"}}
	_, _, err := Cast(entries, opts)
	if err == nil {
		t.Error("expected error casting non-numeric string to int")
	}
}

func TestCast_SkipErrorsContinues(t *testing.T) {
	entries := makeCastEntries()
	opts := CastOptions{
		Types:      map[string]string{"NAME": "int", "PORT": "int"},
		SkipErrors: true,
	}
	_, results, err := Cast(entries, opts)
	if err != nil {
		t.Fatalf("expected no error with SkipErrors=true, got %v", err)
	}
	var hadError bool
	for _, r := range results {
		if r.Error != "" {
			hadError = true
		}
	}
	if !hadError {
		t.Error("expected at least one result with an error message")
	}
}

func TestCast_UnknownTypeReturnsError(t *testing.T) {
	entries := makeCastEntries()
	opts := CastOptions{Types: map[string]string{"PORT": "timestamp"}}
	_, _, err := Cast(entries, opts)
	if err == nil {
		t.Error("expected error for unknown type")
	}
}
