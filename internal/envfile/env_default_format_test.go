package envfile

import (
	"strings"
	"testing"
)

func TestFormatDefaultResults_Empty(t *testing.T) {
	out := FormatDefaultResults(nil, false)
	if !strings.Contains(out, "no defaults") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatDefaultResults_ShowsApplied(t *testing.T) {
	results := []DefaultResult{
		{Key: "LOG_LEVEL", Value: "info", Default: "info", Status: "applied"},
	}
	out := FormatDefaultResults(results, false)
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Error("expected LOG_LEVEL in output")
	}
	if !strings.Contains(out, "+") {
		t.Error("expected '+' marker for applied")
	}
}

func TestFormatDefaultResults_ShowsSkipped(t *testing.T) {
	results := []DefaultResult{
		{Key: "APP_ENV", Value: "production", Default: "staging", Status: "skipped"},
	}
	out := FormatDefaultResults(results, false)
	if !strings.Contains(out, "skipped") {
		t.Error("expected 'skipped' in output")
	}
}

func TestFormatDefaultResults_RedactsSensitive(t *testing.T) {
	results := []DefaultResult{
		{Key: "SECRET_TOKEN", Value: "mysecret", Default: "mysecret", Status: "applied"},
	}
	out := FormatDefaultResults(results, false)
	if strings.Contains(out, "mysecret") {
		t.Error("sensitive value should be redacted")
	}
	if !strings.Contains(out, "[redacted]") {
		t.Error("expected [redacted] in output")
	}
}

func TestFormatDefaultResults_ColorizeAddsEscapeCodes(t *testing.T) {
	results := []DefaultResult{
		{Key: "LOG_LEVEL", Value: "info", Default: "info", Status: "applied"},
	}
	out := FormatDefaultResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI escape codes when colorize=true")
	}
}

func TestFormatDefaultSummary_Counts(t *testing.T) {
	sum := DefaultSummary{Applied: 3, Skipped: 1, Unchanged: 2}
	out := FormatDefaultSummary(sum)
	if !strings.Contains(out, "3 applied") {
		t.Errorf("expected '3 applied' in summary, got: %s", out)
	}
	if !strings.Contains(out, "1 skipped") {
		t.Errorf("expected '1 skipped' in summary, got: %s", out)
	}
}
