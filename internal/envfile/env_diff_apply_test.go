package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func makeDiffApplyBase() []Entry {
	return []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "5432"},
		{Key: "SECRET_KEY", Value: "topsecret"},
	}
}

func makeDiffApplyChanges() []DiffEntry {
	return []DiffEntry{
		{Key: "HOST", Status: "changed", OldValue: "localhost", NewValue: "prod.example.com"},
		{Key: "DEBUG", Status: "added", NewValue: "false"},
		{Key: "PORT", Status: "removed", OldValue: "5432"},
	}
}

func TestApplyDiff_AddsNewKey(t *testing.T) {
	base := makeDiffApplyBase()
	changes := []DiffEntry{{Key: "NEW_VAR", Status: "added", NewValue: "hello"}}
	out, results := ApplyDiff(base, changes, DiffApplyOptions{Overwrite: true})
	if len(results) != 1 || results[0].Action != "added" {
		t.Fatalf("expected added result, got %+v", results)
	}
	m := ToMap(out)
	if m["NEW_VAR"] != "hello" {
		t.Errorf("expected NEW_VAR=hello, got %s", m["NEW_VAR"])
	}
}

func TestApplyDiff_RemovesKey(t *testing.T) {
	base := makeDiffApplyBase()
	changes := []DiffEntry{{Key: "PORT", Status: "removed", OldValue: "5432"}}
	out, results := ApplyDiff(base, changes, DiffApplyOptions{})
	if results[0].Action != "removed" {
		t.Errorf("expected removed, got %s", results[0].Action)
	}
	if _, ok := ToMap(out)["PORT"]; ok {
		t.Error("PORT should have been removed")
	}
}

func TestApplyDiff_UpdatesWithOverwrite(t *testing.T) {
	base := makeDiffApplyBase()
	changes := []DiffEntry{{Key: "HOST", Status: "changed", OldValue: "localhost", NewValue: "prod"}}
	out, results := ApplyDiff(base, changes, DiffApplyOptions{Overwrite: true})
	if results[0].Action != "updated" {
		t.Errorf("expected updated, got %s", results[0].Action)
	}
	if ToMap(out)["HOST"] != "prod" {
		t.Errorf("expected prod, got %s", ToMap(out)["HOST"])
	}
}

func TestApplyDiff_SkipsChangeWithoutOverwrite(t *testing.T) {
	base := makeDiffApplyBase()
	changes := []DiffEntry{{Key: "HOST", Status: "changed", OldValue: "localhost", NewValue: "prod"}}
	_, results := ApplyDiff(base, changes, DiffApplyOptions{Overwrite: false})
	if results[0].Action != "skipped" {
		t.Errorf("expected skipped, got %s", results[0].Action)
	}
}

func TestApplyDiff_SkipsSensitiveKey(t *testing.T) {
	base := makeDiffApplyBase()
	changes := []DiffEntry{{Key: "SECRET_KEY", Status: "changed", OldValue: "topsecret", NewValue: "newsecret"}}
	_, results := ApplyDiff(base, changes, DiffApplyOptions{Overwrite: true, SkipSensitive: true})
	if results[0].Action != "skipped" {
		t.Errorf("expected skipped for sensitive key, got %s", results[0].Action)
	}
}

func TestApplyDiff_DryRunDoesNotMutate(t *testing.T) {
	base := makeDiffApplyBase()
	origLen := len(base)
	changes := []DiffEntry{{Key: "NEWKEY", Status: "added", NewValue: "val"}}
	out, _ := ApplyDiff(base, changes, DiffApplyOptions{DryRun: true})
	if len(out) != origLen {
		t.Errorf("dry run should not add entries, got len %d", len(out))
	}
}

func TestSaveAndLoadDiffApplyResults_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "results.json")
	results := []DiffApplyResult{
		{Key: "A", Action: "added", NewValue: "1"},
		{Key: "B", Action: "removed", OldValue: "old"},
	}
	if err := SaveDiffApplyResults(path, results); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadDiffApplyResults(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 2 || loaded[0].Key != "A" {
		t.Errorf("unexpected loaded results: %+v", loaded)
	}
}

func TestLoadDiffApplyResults_FileNotFound_ReturnsEmpty(t *testing.T) {
	results, err := LoadDiffApplyResults("/nonexistent/path/results.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty, got %d", len(results))
	}
}

func TestSaveDiffApplyResults_RestrictedPerms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "results.json")
	if err := SaveDiffApplyResults(path, []DiffApplyResult{}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 perms, got %v", info.Mode().Perm())
	}
}
