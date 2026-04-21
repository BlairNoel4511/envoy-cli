package envfile

import (
	"strings"
	"testing"
)

func makeSortEntries() []Entry {
	return []Entry{
		{Key: "ZEBRA", Value: "z"},
		{Key: "APPLE", Value: "a"},
		{Key: "MANGO", Value: "m"},
		{Key: "BANANA", Value: "b"},
	}
}

func TestSort_AscendingByKey(t *testing.T) {
	entries := makeSortEntries()
	r := Sort(entries, SortOptions{Order: SortAsc})
	keys := make([]string, len(r.Entries))
	for i, e := range r.Entries {
		keys[i] = e.Key
	}
	expected := []string{"APPLE", "BANANA", "MANGO", "ZEBRA"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("pos %d: want %s, got %s", i, k, keys[i])
		}
	}
}

func TestSort_DescendingByKey(t *testing.T) {
	entries := makeSortEntries()
	r := Sort(entries, SortOptions{Order: SortDesc})
	if r.Entries[0].Key != "ZEBRA" {
		t.Errorf("expected ZEBRA first, got %s", r.Entries[0].Key)
	}
	if r.Entries[len(r.Entries)-1].Key != "APPLE" {
		t.Errorf("expected APPLE last, got %s", r.Entries[len(r.Entries)-1].Key)
	}
}

func TestSort_ByValue(t *testing.T) {
	entries := makeSortEntries()
	r := Sort(entries, SortOptions{Order: SortAsc, ByValue: true})
	if r.Entries[0].Value != "a" {
		t.Errorf("expected value 'a' first, got %s", r.Entries[0].Value)
	}
}

func TestSort_AlreadySorted_ReorderedZero(t *testing.T) {
	entries := []Entry{
		{Key: "ALPHA", Value: "1"},
		{Key: "BETA", Value: "2"},
		{Key: "GAMMA", Value: "3"},
	}
	r := Sort(entries, SortOptions{Order: SortAsc})
	if r.Reordered != 0 {
		t.Errorf("expected 0 reordered, got %d", r.Reordered)
	}
}

func TestSort_OriginalUnmodified(t *testing.T) {
	entries := makeSortEntries()
	origFirst := entries[0].Key
	Sort(entries, SortOptions{Order: SortAsc})
	if entries[0].Key != origFirst {
		t.Error("original slice should not be modified")
	}
}

func TestFormatSortResults_ContainsKeys(t *testing.T) {
	entries := makeSortEntries()
	r := Sort(entries, SortOptions{Order: SortAsc})
	out := FormatSortResults(r, false)
	for _, e := range r.Entries {
		if !strings.Contains(out, e.Key) {
			t.Errorf("expected output to contain key %s", e.Key)
		}
	}
}

func TestFormatSortResults_ColorizeAddsEscapeCodes(t *testing.T) {
	entries := makeSortEntries()
	r := Sort(entries, SortOptions{Order: SortAsc})
	out := FormatSortResults(r, true)
	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI escape codes in colorized output")
	}
}

func TestFormatSortSummaryLine_NoChanges(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	r := Sort(entries, SortOptions{Order: SortAsc})
	summary := FormatSortSummaryLine(r, false)
	if !strings.Contains(summary, "already sorted") {
		t.Errorf("unexpected summary: %s", summary)
	}
}
