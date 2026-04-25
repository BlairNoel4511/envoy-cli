package envfile

import (
	"strings"
	"testing"
)

func makeInspectEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DB_PASSWORD", Value: "s3cr3t", Comment: "database password"},
		{Key: "PORT", Value: "8080"},
		{Key: "API_SECRET", Value: "topsecret"},
	}
}

func TestInspect_AllEntries(t *testing.T) {
	entries := makeInspectEntries()
	results := Inspect(entries, nil, InspectOptions{Redact: false})
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
}

func TestInspect_SpecificKeys(t *testing.T) {
	entries := makeInspectEntries()
	results := Inspect(entries, []string{"PORT", "APP_NAME"}, InspectOptions{})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Found {
			t.Errorf("expected key %q to be found", r.Key)
		}
	}
}

func TestInspect_MissingKeyReturnsNotFound(t *testing.T) {
	entries := makeInspectEntries()
	results := Inspect(entries, []string{"MISSING_KEY"}, InspectOptions{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Found {
		t.Error("expected Found=false for missing key")
	}
}

func TestInspect_RedactsSensitiveValues(t *testing.T) {
	entries := makeInspectEntries()
	results := Inspect(entries, []string{"DB_PASSWORD", "API_SECRET"}, InspectOptions{Redact: true})
	for _, r := range results {
		if r.Value != "***" {
			t.Errorf("expected redacted value for %q, got %q", r.Key, r.Value)
		}
		if !r.Redacted {
			t.Errorf("expected Redacted=true for %q", r.Key)
		}
	}
}

func TestInspect_DetectsComment(t *testing.T) {
	entries := makeInspectEntries()
	results := Inspect(entries, []string{"DB_PASSWORD"}, InspectOptions{})
	if len(results) == 0 || !results[0].HasComment {
		t.Error("expected HasComment=true for DB_PASSWORD")
	}
	if results[0].Comment != "database password" {
		t.Errorf("unexpected comment: %q", results[0].Comment)
	}
}

func TestInspect_LengthIsCorrect(t *testing.T) {
	entries := makeInspectEntries()
	results := Inspect(entries, []string{"APP_NAME"}, InspectOptions{})
	if results[0].Length != len("envoy") {
		t.Errorf("expected length %d, got %d", len("envoy"), results[0].Length)
	}
}

func TestFormatInspectResults_NotFound(t *testing.T) {
	results := []InspectResult{{Key: "GHOST", Found: false}}
	out := FormatInspectResults(results, false)
	if !strings.Contains(out, "not found") {
		t.Errorf("expected 'not found' in output, got: %q", out)
	}
}

func TestFormatInspectResults_ColorizeAddsEscapeCodes(t *testing.T) {
	results := []InspectResult{
		{Key: "API_SECRET", Value: "***", Found: true, Sensitive: true, Redacted: true, Length: 9},
	}
	out := FormatInspectResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes in colorized output")
	}
}

func TestInspectSummary_Counts(t *testing.T) {
	results := []InspectResult{
		{Key: "A", Found: true, Sensitive: false},
		{Key: "B", Found: true, Sensitive: true},
		{Key: "C", Found: false},
	}
	summary := InspectSummary(results)
	if !strings.Contains(summary, "3 inspected") {
		t.Errorf("unexpected summary: %q", summary)
	}
	if !strings.Contains(summary, "2 found") {
		t.Errorf("unexpected summary: %q", summary)
	}
	if !strings.Contains(summary, "1 sensitive") {
		t.Errorf("unexpected summary: %q", summary)
	}
}
