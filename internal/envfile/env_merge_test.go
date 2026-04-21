package envfile

import (
	"testing"
)

func makeMergeEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestMergeMany_AddsNewKeys(t *testing.T) {
	dst := makeMergeEntries("APP", "1")
	src := makeMergeEntries("DB", "postgres")
	out, results := MergeMany(dst, src, MergeOptions{})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if results[0].Action != "added" {
		t.Errorf("expected added, got %s", results[0].Action)
	}
}

func TestMergeMany_SkipsExistingWithoutOverwrite(t *testing.T) {
	dst := makeMergeEntries("APP", "1")
	src := makeMergeEntries("APP", "2")
	out, results := MergeMany(dst, src, MergeOptions{Overwrite: false})
	if out[0].Value != "1" {
		t.Errorf("expected original value preserved")
	}
	if results[0].Action != "skipped" {
		t.Errorf("expected skipped, got %s", results[0].Action)
	}
}

func TestMergeMany_OverwritesWhenEnabled(t *testing.T) {
	dst := makeMergeEntries("APP", "1")
	src := makeMergeEntries("APP", "2")
	out, results := MergeMany(dst, src, MergeOptions{Overwrite: true})
	if out[0].Value != "2" {
		t.Errorf("expected updated value")
	}
	if results[0].Action != "updated" {
		t.Errorf("expected updated, got %s", results[0].Action)
	}
}

func TestMergeMany_SkipsSensitiveKeys(t *testing.T) {
	dst := makeMergeEntries()
	src := makeMergeEntries("SECRET_TOKEN", "abc")
	_, results := MergeMany(dst, src, MergeOptions{SkipSensitive: true})
	if results[0].Action != "skipped" {
		t.Errorf("expected sensitive key to be skipped, got %s", results[0].Action)
	}
}

func TestMergeMany_DryRunDoesNotWrite(t *testing.T) {
	dst := makeMergeEntries("APP", "1")
	src := makeMergeEntries("NEW", "val")
	out, results := MergeMany(dst, src, MergeOptions{DryRun: true})
	if len(out) != 1 {
		t.Errorf("dry run should not modify dst, got %d entries", len(out))
	}
	if results[0].Action != "added" {
		t.Errorf("expected added result even in dry run")
	}
}

func TestMergeMany_UnchangedWhenValueIdentical(t *testing.T) {
	dst := makeMergeEntries("APP", "same")
	src := makeMergeEntries("APP", "same")
	_, results := MergeMany(dst, src, MergeOptions{Overwrite: true})
	if results[0].Action != "unchanged" {
		t.Errorf("expected unchanged, got %s", results[0].Action)
	}
}

func TestMergeManySum_Counts(t *testing.T) {
	results := []MergeResult{
		{Action: "added"}, {Action: "added"},
		{Action: "updated"},
		{Action: "skipped"},
		{Action: "unchanged"}, {Action: "unchanged"},
	}
	s := MergeManySum(results)
	if s.Added != 2 || s.Updated != 1 || s.Skipped != 1 || s.Unchanged != 2 {
		t.Errorf("unexpected summary: %+v", s)
	}
}
