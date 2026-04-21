package envfile

import (
	"strings"
	"testing"
)

func TestFormatDuplicateResult_Duplicated(t *testing.T) {
	r := DuplicateResult{SourceKey: "A", DestKey: "B", Status: "duplicated", Value: "hello"}
	out := FormatDuplicateResult(r, false)
	if !strings.Contains(out, "A => B") {
		t.Errorf("expected arrow notation, got %q", out)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("expected value in output, got %q", out)
	}
}

func TestFormatDuplicateResult_RedactsSensitive(t *testing.T) {
	r := DuplicateResult{SourceKey: "SECRET_KEY", DestKey: "SECRET_COPY", Status: "duplicated", Value: "s3cr3t", Sensitive: true}
	out := FormatDuplicateResult(r, false)
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("expected value to be redacted, got %q", out)
	}
	if !strings.Contains(out, "[redacted]") {
		t.Errorf("expected [redacted] in output, got %q", out)
	}
}

func TestFormatDuplicateResult_SourceNotFound(t *testing.T) {
	r := DuplicateResult{SourceKey: "MISSING", DestKey: "DEST", Status: "source_not_found"}
	out := FormatDuplicateResult(r, false)
	if !strings.Contains(out, "source not found") {
		t.Errorf("expected 'source not found' in output, got %q", out)
	}
}

func TestFormatDuplicateResult_ColorizeAddsEscapeCodes(t *testing.T) {
	r := DuplicateResult{SourceKey: "A", DestKey: "B", Status: "duplicated", Value: "v"}
	out := FormatDuplicateResult(r, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes in colorized output, got %q", out)
	}
}

func TestFormatDuplicateSummary_Counts(t *testing.T) {
	s := DuplicateSummary{Duplicated: 3, Skipped: 1, Unchanged: 2, NotFound: 0}
	out := FormatDuplicateSummary(s)
	if !strings.Contains(out, "3 duplicated") {
		t.Errorf("expected '3 duplicated' in summary, got %q", out)
	}
	if !strings.Contains(out, "1 skipped") {
		t.Errorf("expected '1 skipped' in summary, got %q", out)
	}
	if !strings.Contains(out, "2 unchanged") {
		t.Errorf("expected '2 unchanged' in summary, got %q", out)
	}
}
