package envfile

import (
	"strings"
	"testing"
)

func makeCompareEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestCompare_OnlyInLeft(t *testing.T) {
	left := makeCompareEntries("FOO", "bar", "BAZ", "qux")
	right := makeCompareEntries("FOO", "bar")
	r := Compare(left, right)
	if len(r.OnlyInLeft) != 1 || r.OnlyInLeft[0].Key != "BAZ" {
		t.Errorf("expected BAZ in OnlyInLeft, got %+v", r.OnlyInLeft)
	}
}

func TestCompare_OnlyInRight(t *testing.T) {
	left := makeCompareEntries("FOO", "bar")
	right := makeCompareEntries("FOO", "bar", "NEW", "val")
	r := Compare(left, right)
	if len(r.OnlyInRight) != 1 || r.OnlyInRight[0].Key != "NEW" {
		t.Errorf("expected NEW in OnlyInRight, got %+v", r.OnlyInRight)
	}
}

func TestCompare_Different(t *testing.T) {
	left := makeCompareEntries("FOO", "old")
	right := makeCompareEntries("FOO", "new")
	r := Compare(left, right)
	if len(r.Different) != 1 || r.Different[0].LeftValue != "old" || r.Different[0].RightValue != "new" {
		t.Errorf("unexpected Different: %+v", r.Different)
	}
}

func TestCompare_Identical(t *testing.T) {
	left := makeCompareEntries("FOO", "bar")
	right := makeCompareEntries("FOO", "bar")
	r := Compare(left, right)
	if len(r.Identical) != 1 || r.Identical[0].Key != "FOO" {
		t.Errorf("expected FOO in Identical, got %+v", r.Identical)
	}
}

func TestHasDifferences_True(t *testing.T) {
	left := makeCompareEntries("A", "1")
	right := makeCompareEntries("B", "2")
	r := Compare(left, right)
	if !HasDifferences(r) {
		t.Error("expected HasDifferences to be true")
	}
}

func TestHasDifferences_False(t *testing.T) {
	left := makeCompareEntries("A", "1")
	right := makeCompareEntries("A", "1")
	r := Compare(left, right)
	if HasDifferences(r) {
		t.Error("expected HasDifferences to be false")
	}
}

func TestFormatCompareResult_ContainsSymbols(t *testing.T) {
	left := makeCompareEntries("OLD", "x")
	right := makeCompareEntries("NEW", "y")
	r := Compare(left, right)
	out := FormatCompareResult(r, false)
	if !strings.Contains(out, "+ NEW") {
		t.Errorf("expected '+ NEW' in output: %s", out)
	}
	if !strings.Contains(out, "- OLD") {
		t.Errorf("expected '- OLD' in output: %s", out)
	}
}

func TestFormatCompareResult_RedactsSensitive(t *testing.T) {
	left := makeCompareEntries("SECRET_KEY", "mysecret")
	right := makeCompareEntries("SECRET_KEY", "newsecret")
	r := Compare(left, right)
	out := FormatCompareResult(r, true)
	if strings.Contains(out, "mysecret") || strings.Contains(out, "newsecret") {
		t.Errorf("expected sensitive values to be redacted: %s", out)
	}
}

func TestFormatCompareSummaryLine_NoDiff(t *testing.T) {
	left := makeCompareEntries("A", "1")
	right := makeCompareEntries("A", "1")
	r := Compare(left, right)
	line := FormatCompareSummaryLine(r)
	if line != "No differences found." {
		t.Errorf("unexpected summary: %s", line)
	}
}
