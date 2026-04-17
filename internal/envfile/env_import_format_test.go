package envfile

import (
	"strings"
	"testing"
)

func TestFormatImportResults_Empty(t *testing.T) {
	out := FormatImportResults(nil, false)
	if !strings.Contains(out, "nothing") {
		t.Errorf("expected empty message, got %q", out)
	}
}

func TestFormatImportResults_ShowsAdded(t *testing.T) {
	results := []ImportResult{{Key: "FOO", Value: "bar", Status: "added"}}
	out := FormatImportResults(results, false)
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "added") {
		t.Errorf("expected FOO added in output, got %q", out)
	}
}

func TestFormatImportResults_RedactsSensitive(t *testing.T) {
	results := []ImportResult{{Key: "SECRET", Value: "topsecret", Status: "added", Sensitive: true}}
	out := FormatImportResults(results, false)
	if strings.Contains(out, "topsecret") {
		t.Errorf("sensitive value should be redacted")
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected *** in output")
	}
}

func TestFormatImportResults_SortedOutput(t *testing.T) {
	results := []ImportResult{
		{Key: "Z_KEY", Value: "z", Status: "added"},
		{Key: "A_KEY", Value: "a", Status: "added"},
	}
	out := FormatImportResults(results, false)
	idxA := strings.Index(out, "A_KEY")
	idxZ := strings.Index(out, "Z_KEY")
	if idxA > idxZ {
		t.Errorf("expected A_KEY before Z_KEY in output")
	}
}

func TestFormatImportResults_ColorizeAddsEscapeCodes(t *testing.T) {
	results := []ImportResult{{Key: "FOO", Value: "bar", Status: "added"}}
	out := FormatImportResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes in colorized output")
	}
}

func TestFormatImportSummary_Counts(t *testing.T) {
	s := ImportSummary{Added: 3, Overwritten: 1, Skipped: 2, Total: 6}
	out := FormatImportSummary(s)
	if !strings.Contains(out, "3 added") || !strings.Contains(out, "1 overwritten") || !strings.Contains(out, "2 skipped") {
		t.Errorf("unexpected summary output: %q", out)
	}
}
