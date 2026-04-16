package envfile

import (
	"strings"
	"testing"
)

func makeDeleteEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: "secret"},
	}
}

func TestDelete_RemovesExistingKey(t *testing.T) {
	entries := makeDeleteEntries()
	updated, res := Delete(entries, "DB_HOST")
	if !res.Deleted {
		t.Fatalf("expected key to be deleted")
	}
	for _, e := range updated {
		if e.Key == "DB_HOST" {
			t.Errorf("expected DB_HOST to be removed from entries")
		}
	}
	if len(updated) != 2 {
		t.Errorf("expected 2 entries, got %d", len(updated))
	}
}

func TestDelete_KeyNotFound(t *testing.T) {
	entries := makeDeleteEntries()
	_, res := Delete(entries, "MISSING_KEY")
	if res.Deleted {
		t.Errorf("expected Deleted=false for missing key")
	}
	if res.Reason == "" {
		t.Errorf("expected a reason for skipped deletion")
	}
}

func TestDeleteMany_DeletesMultipleKeys(t *testing.T) {
	entries := makeDeleteEntries()
	updated, results, summary := DeleteMany(entries, []string{"APP_NAME", "DB_HOST"})
	if summary.Deleted != 2 {
		t.Errorf("expected 2 deleted, got %d", summary.Deleted)
	}
	if summary.Skipped != 0 {
		t.Errorf("expected 0 skipped, got %d", summary.Skipped)
	}
	if len(updated) != 1 {
		t.Errorf("expected 1 remaining entry, got %d", len(updated))
	}
	_ = results
}

func TestDeleteMany_SkipsMissingKeys(t *testing.T) {
	entries := makeDeleteEntries()
	_, _, summary := DeleteMany(entries, []string{"NONEXISTENT"})
	if summary.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", summary.Skipped)
	}
}

func TestFormatDeleteResults_ShowsDeleted(t *testing.T) {
	results := []DeleteResult{
		{Key: "DB_HOST", Deleted: true, Reason: "deleted"},
	}
	out := FormatDeleteResults(results, false)
	if !strings.Contains(out, "- DB_HOST") {
		t.Errorf("expected deleted key in output, got: %s", out)
	}
}

func TestFormatDeleteResults_ShowsSkipped(t *testing.T) {
	results := []DeleteResult{
		{Key: "MISSING", Deleted: false, Reason: "key not found"},
	}
	out := FormatDeleteResults(results, false)
	if !strings.Contains(out, "~ MISSING") {
		t.Errorf("expected skipped key in output, got: %s", out)
	}
}

func TestFormatDeleteResults_ColorizeAddsEscapeCodes(t *testing.T) {
	results := []DeleteResult{
		{Key: "APP_NAME", Deleted: true, Reason: "deleted"},
	}
	out := FormatDeleteResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes in colorized output")
	}
}

func TestFormatDeleteSummary_Counts(t *testing.T) {
	s := DeleteSummary{Deleted: 3, Skipped: 1}
	out := FormatDeleteSummary(s)
	if !strings.Contains(out, "3 deleted") || !strings.Contains(out, "1 skipped") {
		t.Errorf("unexpected summary format: %s", out)
	}
}
