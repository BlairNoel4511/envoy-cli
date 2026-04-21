package envfile

import (
	"strings"
	"testing"
)

func TestFormatMergeResults_Empty(t *testing.T) {
	out := FormatMergeResults(nil, false)
	if out != "(no changes)" {
		t.Errorf("expected no-changes message, got %q", out)
	}
}

func TestFormatMergeResults_ShowsAdded(t *testing.T) {
	results := []MergeResult{{Key: "FOO", Value: "bar", Action: "added"}}
	out := FormatMergeResults(results, false)
	if !strings.Contains(out, "+ FOO=bar") {
		t.Errorf("expected added line, got %q", out)
	}
}

func TestFormatMergeResults_ShowsUpdated(t *testing.T) {
	results := []MergeResult{{Key: "FOO", Value: "new", Action: "updated"}}
	out := FormatMergeResults(results, false)
	if !strings.Contains(out, "~ FOO=new") {
		t.Errorf("expected updated line, got %q", out)
	}
}

func TestFormatMergeResults_ShowsSkipped(t *testing.T) {
	results := []MergeResult{{Key: "FOO", Value: "v", Action: "skipped"}}
	out := FormatMergeResults(results, false)
	if !strings.Contains(out, "! FOO=v") {
		t.Errorf("expected skipped line, got %q", out)
	}
}

func TestFormatMergeResults_RedactsSensitive(t *testing.T) {
	results := []MergeResult{{Key: "SECRET_KEY", Value: "topsecret", Action: "added", Sensitive: true}}
	out := FormatMergeResults(results, false)
	if strings.Contains(out, "topsecret") {
		t.Errorf("sensitive value should be redacted")
	}
	if !strings.Contains(out, "[redacted]") {
		t.Errorf("expected redacted placeholder")
	}
}

func TestFormatMergeResults_SortedOutput(t *testing.T) {
	results := []MergeResult{
		{Key: "ZZZ", Value: "z", Action: "added"},
		{Key: "AAA", Value: "a", Action: "added"},
	}
	out := FormatMergeResults(results, false)
	lines := strings.Split(out, "\n")
	if !strings.Contains(lines[0], "AAA") {
		t.Errorf("expected AAA first, got %q", lines[0])
	}
}

func TestFormatMergeResults_ColorizeAddsEscapeCodes(t *testing.T) {
	results := []MergeResult{{Key: "FOO", Value: "bar", Action: "added"}}
	out := FormatMergeResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes with colorize=true")
	}
}

func TestFormatMergeSummaryLine_Counts(t *testing.T) {
	s := MergeSummary{Added: 3, Updated: 1, Skipped: 2, Unchanged: 5}
	line := FormatMergeSummaryLine(s)
	if !strings.Contains(line, "3 added") || !strings.Contains(line, "1 updated") {
		t.Errorf("unexpected summary line: %q", line)
	}
}
