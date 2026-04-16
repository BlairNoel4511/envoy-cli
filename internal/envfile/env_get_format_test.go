package envfile

import (
	"strings"
	"testing"
)

func TestFormatGetResults_Empty(t *testing.T) {
	out := FormatGetResults(nil, false)
	if out != "(no keys requested)" {
		t.Errorf("unexpected: %s", out)
	}
}

func TestFormatGetResults_ShowsFound(t *testing.T) {
	results := []GetResult{
		{Key: "PORT", Value: "8080", Found: true},
	}
	out := FormatGetResults(results, false)
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT=8080 in output, got: %s", out)
	}
}

func TestFormatGetResults_ShowsNotFound(t *testing.T) {
	results := []GetResult{
		{Key: "MISSING", Found: false},
	}
	out := FormatGetResults(results, false)
	if !strings.Contains(out, "not found") {
		t.Errorf("expected 'not found' in output, got: %s", out)
	}
}

func TestFormatGetResults_ShowsRedacted(t *testing.T) {
	results := []GetResult{
		{Key: "DB_PASSWORD", Value: "***", Found: true, Redacted: true},
	}
	out := FormatGetResults(results, false)
	if !strings.Contains(out, "***") {
		t.Errorf("expected redacted value in output, got: %s", out)
	}
}

func TestFormatGetResults_ColorizeAddsEscapeCodes(t *testing.T) {
	results := []GetResult{
		{Key: "APP", Value: "test", Found: true},
	}
	out := FormatGetResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI codes in colorized output, got: %s", out)
	}
}

func TestFormatGetResults_MultipleLines(t *testing.T) {
	results := []GetResult{
		{Key: "A", Value: "1", Found: true},
		{Key: "B", Found: false},
	}
	out := FormatGetResults(results, false)
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d: %q", len(lines), out)
	}
}
