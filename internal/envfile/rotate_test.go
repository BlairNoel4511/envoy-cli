package envfile

import (
	"strings"
	"testing"
)

func makeRotateEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "old-secret"},
		{Key: "API_KEY", Value: "key-v1"},
	}
}

func TestRotate_UpdatesValues(t *testing.T) {
	entries := makeRotateEntries()
	newVals := map[string]string{"API_KEY": "key-v2"}

	result, summary, err := Rotate(entries, newVals, RotateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ToMap(result)
	if m["API_KEY"] != "key-v2" {
		t.Errorf("expected API_KEY=key-v2, got %s", m["API_KEY"])
	}
	if summary.Rotated != 1 {
		t.Errorf("expected 1 rotated, got %d", summary.Rotated)
	}
}

func TestRotate_SkipsIdenticalValues(t *testing.T) {
	entries := makeRotateEntries()
	newVals := map[string]string{"APP_NAME": "myapp"}

	_, summary, err := Rotate(entries, newVals, RotateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Skipped != 1 || summary.Rotated != 0 {
		t.Errorf("expected 1 skipped, got rotated=%d skipped=%d", summary.Rotated, summary.Skipped)
	}
}

func TestRotate_ErrorOnMissingKeyByDefault(t *testing.T) {
	entries := makeRotateEntries()
	newVals := map[string]string{"MISSING_KEY": "value"}

	_, _, err := Rotate(entries, newVals, RotateOptions{})
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRotate_SkipMissingKey(t *testing.T) {
	entries := makeRotateEntries()
	newVals := map[string]string{"MISSING_KEY": "value"}

	_, summary, err := Rotate(entries, newVals, RotateOptions{SkipMissing: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", summary.Skipped)
	}
}

func TestRotate_DryRunDoesNotMutate(t *testing.T) {
	entries := makeRotateEntries()
	newVals := map[string]string{"API_KEY": "key-v3"}

	_, _, err := Rotate(entries, newVals, RotateOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ToMap(entries)
	if m["API_KEY"] != "key-v1" {
		t.Errorf("dry run should not mutate original entries")
	}
}

func TestFormatRotateSummary_ContainsKey(t *testing.T) {
	summary := RotateSummary{
		Results: []RotateResult{
			{Key: "API_KEY", OldValue: "old", NewValue: "new"},
		},
		Rotated: 1,
	}
	out := FormatRotateSummary(summary, false)
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got: %s", out)
	}
}

func TestFormatRotateSummary_RedactsSensitive(t *testing.T) {
	summary := RotateSummary{
		Results: []RotateResult{
			{Key: "DB_PASSWORD", OldValue: "old-secret", NewValue: "new-secret"},
		},
		Rotated: 1,
	}
	out := FormatRotateSummary(summary, false)
	if strings.Contains(out, "old-secret") || strings.Contains(out, "new-secret") {
		t.Errorf("sensitive values should be redacted, got: %s", out)
	}
	if !strings.Contains(out, "[redacted]") {
		t.Errorf("expected [redacted] in output, got: %s", out)
	}
}
