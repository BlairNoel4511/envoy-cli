package envfile

import (
	"strings"
	"testing"
)

func makeSetEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DEBUG", Value: "false"},
		{Key: "DB_PASSWORD", Value: "secret"},
	}
}

func TestSet_AddsNewKey(t *testing.T) {
	entries := makeSetEntries()
	out, r := Set(entries, "NEW_KEY", "hello", SetOptions{})
	if r.Action != "added" {
		t.Errorf("expected added, got %s", r.Action)
	}
	m := ToMap(out)
	if m["NEW_KEY"] != "hello" {
		t.Errorf("expected hello, got %s", m["NEW_KEY"])
	}
}

func TestSet_SkipsExistingWithoutOverwrite(t *testing.T) {
	entries := makeSetEntries()
	out, r := Set(entries, "APP_NAME", "other", SetOptions{Overwrite: false})
	if r.Action != "skipped" {
		t.Errorf("expected skipped, got %s", r.Action)
	}
	if ToMap(out)["APP_NAME"] != "myapp" {
		t.Error("value should not have changed")
	}
}

func TestSet_UpdatesExistingWithOverwrite(t *testing.T) {
	entries := makeSetEntries()
	out, r := Set(entries, "APP_NAME", "newapp", SetOptions{Overwrite: true})
	if r.Action != "updated" {
		t.Errorf("expected updated, got %s", r.Action)
	}
	if ToMap(out)["APP_NAME"] != "newapp" {
		t.Error("expected value to be updated")
	}
}

func TestSet_UnchangedWhenValueIdentical(t *testing.T) {
	entries := makeSetEntries()
	_, r := Set(entries, "DEBUG", "false", SetOptions{Overwrite: true})
	if r.Action != "unchanged" {
		t.Errorf("expected unchanged, got %s", r.Action)
	}
}

func TestSet_DryRunDoesNotModify(t *testing.T) {
	entries := makeSetEntries()
	out, r := Set(entries, "APP_NAME", "dryval", SetOptions{Overwrite: true, DryRun: true})
	if r.Action != "updated" {
		t.Errorf("expected updated action in dry run, got %s", r.Action)
	}
	if ToMap(out)["APP_NAME"] != "myapp" {
		t.Error("dry run should not modify entries")
	}
}

func TestSetMany_AppliesMultiplePairs(t *testing.T) {
	entries := makeSetEntries()
	pairs := map[string]string{"X": "1", "Y": "2"}
	out, results := SetMany(entries, pairs, SetOptions{})
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	m := ToMap(out)
	if m["X"] != "1" || m["Y"] != "2" {
		t.Error("expected both keys to be set")
	}
}

func TestSetSummary_Counts(t *testing.T) {
	results := []SetResult{
		{Action: "added"}, {Action: "updated"}, {Action: "skipped"}, {Action: "unchanged"},
	}
	summary := SetSummary(results)
	if !strings.Contains(summary, "added=1") || !strings.Contains(summary, "updated=1") {
		t.Errorf("unexpected summary: %s", summary)
	}
}
