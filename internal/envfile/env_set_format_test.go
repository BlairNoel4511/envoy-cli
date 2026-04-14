package envfile

import (
	"strings"
	"testing"
)

func TestFormatSetResults_Empty(t *testing.T) {
	out := FormatSetResults(nil, false)
	if out != "(no changes)" {
		t.Errorf("expected no changes message, got %q", out)
	}
}

func TestFormatSetResults_ShowsAdded(t *testing.T) {
	results := []SetResult{{Key: "FOO", Value: "bar", Action: "added"}}
	out := FormatSetResults(results, false)
	if !strings.Contains(out, "+") || !strings.Contains(out, "FOO") {
		t.Errorf("expected added line, got %q", out)
	}
}

func TestFormatSetResults_ShowsUpdated(t *testing.T) {
	results := []SetResult{{
		Key: "APP", Value: "new", OldValue: "old", Action: "updated",
	}}
	out := FormatSetResults(results, false)
	if !strings.Contains(out, "~") || !strings.Contains(out, "→") {
		t.Errorf("expected updated line with arrow, got %q", out)
	}
}

func TestFormatSetResults_ShowsSkipped(t *testing.T) {
	results := []SetResult{{Key: "X", Value: "v", Action: "skipped"}}
	out := FormatSetResults(results, false)
	if !strings.Contains(out, "!") || !strings.Contains(out, "skipped") {
		t.Errorf("expected skipped line, got %q", out)
	}
}

func TestFormatSetResults_ShowsUnchanged(t *testing.T) {
	results := []SetResult{{Key: "Z", Value: "v", Action: "unchanged"}}
	out := FormatSetResults(results, false)
	if !strings.Contains(out, "=") || !strings.Contains(out, "unchanged") {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestFormatSetResults_ColorizeAddsEscapeCodes(t *testing.T) {
	results := []SetResult{{Key: "FOO", Value: "bar", Action: "added"}}
	out := FormatSetResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes, got %q", out)
	}
}

func TestFormatSetResults_RedactsSensitiveValues(t *testing.T) {
	results := []SetResult{{Key: "DB_PASSWORD", Value: "supersecret", Action: "added"}}
	out := FormatSetResults(results, false)
	if strings.Contains(out, "supersecret") {
		t.Errorf("expected sensitive value to be redacted, got %q", out)
	}
}
