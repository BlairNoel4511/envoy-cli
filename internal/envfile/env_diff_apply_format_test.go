package envfile

import (
	"strings"
	"testing"
)

func TestFormatDiffApplyResults_Empty(t *testing.T) {
	out := FormatDiffApplyResults(nil, false)
	if out != "(no changes)" {
		t.Errorf("expected no changes, got %q", out)
	}
}

func TestFormatDiffApplyResults_ShowsAdded(t *testing.T) {
	results := []DiffApplyResult{{Key: "FOO", Action: "added", NewValue: "bar"}}
	out := FormatDiffApplyResults(results, false)
	if !strings.Contains(out, "+ FOO = bar") {
		t.Errorf("expected added line, got %q", out)
	}
}

func TestFormatDiffApplyResults_ShowsUpdated(t *testing.T) {
	results := []DiffApplyResult{{
		Key: "HOST", Action: "updated", OldValue: "localhost", NewValue: "prod",
	}}
	out := FormatDiffApplyResults(results, false)
	if !strings.Contains(out, "~ HOST") || !strings.Contains(out, "localhost") || !strings.Contains(out, "prod") {
		t.Errorf("expected updated line with old and new, got %q", out)
	}
}

func TestFormatDiffApplyResults_ShowsRemoved(t *testing.T) {
	results := []DiffApplyResult{{Key: "PORT", Action: "removed", OldValue: "5432"}}
	out := FormatDiffApplyResults(results, false)
	if !strings.Contains(out, "- PORT") {
		t.Errorf("expected removed line, got %q", out)
	}
}

func TestFormatDiffApplyResults_RedactsSensitive(t *testing.T) {
	results := []DiffApplyResult{
		{Key: "SECRET_KEY", Action: "added", NewValue: "mysecret", Sensitive: true},
	}
	out := FormatDiffApplyResults(results, false)
	if strings.Contains(out, "mysecret") {
		t.Errorf("sensitive value should be redacted, got %q", out)
	}
	if !strings.Contains(out, "[redacted]") {
		t.Errorf("expected [redacted], got %q", out)
	}
}

func TestFormatDiffApplyResults_ColorizeAddsEscapeCodes(t *testing.T) {
	results := []DiffApplyResult{{Key: "A", Action: "added", NewValue: "1"}}
	out := FormatDiffApplyResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes in colorized output, got %q", out)
	}
}

func TestFormatDiffApplySummary_Counts(t *testing.T) {
	s := DiffApplySummary{Added: 2, Updated: 1, Removed: 3, Skipped: 0}
	out := FormatDiffApplySummary(s)
	if !strings.Contains(out, "+2") || !strings.Contains(out, "~1") || !strings.Contains(out, "-3") {
		t.Errorf("unexpected summary output: %q", out)
	}
}
