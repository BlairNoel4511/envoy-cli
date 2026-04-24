package envfile

import (
	"strings"
	"testing"
)

func makeExistsEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
		{Key: "DEBUG", Value: "true"},
	}
}

func TestExists_KeyFound(t *testing.T) {
	entries := makeExistsEntries()
	r := Exists(entries, "APP_NAME", ExistsOptions{})
	if !r.Exists {
		t.Fatal("expected key to exist")
	}
	if r.Value != "envoy" {
		t.Errorf("unexpected value: %s", r.Value)
	}
}

func TestExists_KeyNotFound(t *testing.T) {
	entries := makeExistsEntries()
	r := Exists(entries, "MISSING", ExistsOptions{})
	if r.Exists {
		t.Fatal("expected key to be missing")
	}
}

func TestExists_RedactsSensitiveKey(t *testing.T) {
	entries := makeExistsEntries()
	r := Exists(entries, "SECRET_KEY", ExistsOptions{RedactSensitive: true})
	if !r.Exists {
		t.Fatal("expected key to exist")
	}
	if r.Value != "***" {
		t.Errorf("expected redacted value, got %s", r.Value)
	}
	if !r.Masked {
		t.Error("expected Masked to be true")
	}
}

func TestExists_NoRedactWhenDisabled(t *testing.T) {
	entries := makeExistsEntries()
	r := Exists(entries, "SECRET_KEY", ExistsOptions{RedactSensitive: false})
	if r.Value == "***" {
		t.Error("value should not be redacted")
	}
}

func TestExistsMany_MixedResults(t *testing.T) {
	entries := makeExistsEntries()
	results := ExistsMany(entries, []string{"APP_NAME", "MISSING", "DEBUG"}, ExistsOptions{})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if !results[0].Exists || results[1].Exists || !results[2].Exists {
		t.Error("unexpected exists pattern")
	}
}

func TestSummarizeExists_Counts(t *testing.T) {
	results := []ExistsResult{
		{Key: "A", Exists: true},
		{Key: "B", Exists: false},
		{Key: "C", Exists: true},
	}
	s := SummarizeExists(results)
	if s.Found != 2 || s.Missing != 1 {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestFormatExistsResults_ShowsFound(t *testing.T) {
	results := []ExistsResult{{Key: "APP_NAME", Exists: true, Value: "envoy"}}
	out := FormatExistsResults(results, false)
	if !strings.Contains(out, "APP_NAME") {
		t.Error("expected key in output")
	}
	if !strings.Contains(out, "✔") {
		t.Error("expected checkmark in output")
	}
}

func TestFormatExistsResults_ShowsNotFound(t *testing.T) {
	results := []ExistsResult{{Key: "MISSING", Exists: false}}
	out := FormatExistsResults(results, false)
	if !strings.Contains(out, "✘") {
		t.Error("expected cross in output")
	}
}

func TestFormatExistsResults_ColorizeAddsEscapeCodes(t *testing.T) {
	results := []ExistsResult{{Key: "APP_NAME", Exists: true, Value: "v"}}
	out := FormatExistsResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI escape codes in colorized output")
	}
}
