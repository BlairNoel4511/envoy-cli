package envfile

import (
	"strings"
	"testing"
)

func TestFormatCompareSummaryLine_WithDiffs(t *testing.T) {
	left := makeCompareEntries("OLD", "v1", "SHARED", "v")
	right := makeCompareEntries("NEW", "v2", "SHARED", "v2")
	r := Compare(left, right)
	line := FormatCompareSummaryLine(r)
	if !strings.Contains(line, "added") {
		t.Errorf("expected 'added' in summary: %s", line)
	}
	if !strings.Contains(line, "removed") {
		t.Errorf("expected 'removed' in summary: %s", line)
	}
	if !strings.Contains(line, "changed") {
		t.Errorf("expected 'changed' in summary: %s", line)
	}
}

func TestFormatCompareResult_ShowsIdentical(t *testing.T) {
	left := makeCompareEntries("SAME", "val")
	right := makeCompareEntries("SAME", "val")
	r := Compare(left, right)
	out := FormatCompareResult(r, false)
	if !strings.Contains(out, "  SAME=val") {
		t.Errorf("expected identical entry with leading spaces: %s", out)
	}
}

func TestFormatCompareResult_ShowsChange(t *testing.T) {
	left := makeCompareEntries("PORT", "3000")
	right := makeCompareEntries("PORT", "8080")
	r := Compare(left, right)
	out := FormatCompareResult(r, false)
	if !strings.Contains(out, "~ PORT: 3000 -> 8080") {
		t.Errorf("expected change line: %s", out)
	}
}

func TestCompareSummary_Format(t *testing.T) {
	left := makeCompareEntries("A", "1", "B", "2")
	right := makeCompareEntries("A", "1", "C", "3")
	r := Compare(left, right)
	s := CompareSummary(r)
	if !strings.Contains(s, "added") || !strings.Contains(s, "removed") {
		t.Errorf("unexpected summary format: %s", s)
	}
}

func TestFormatCompareResult_EmptyBothSides(t *testing.T) {
	r := Compare([]Entry{}, []Entry{})
	out := FormatCompareResult(r, false)
	if out != "" {
		t.Errorf("expected empty output for empty inputs, got: %q", out)
	}
}
