package envfile

import (
	"strings"
	"testing"
)

func makeRenameMany() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_PASS", Value: "secret"},
	}
}

func TestRenameMany_BasicRename(t *testing.T) {
	entries := makeRenameMany()
	out, results := RenameMany(entries, map[string]string{"APP_HOST": "HOST"}, RenameManyOptions{})
	_, found := Lookup(out, "HOST")
	if !found {
		t.Fatal("expected HOST to exist after rename")
	}
	_, old := Lookup(out, "APP_HOST")
	if old {
		t.Fatal("expected APP_HOST to be gone")
	}
	if results[0].Status != "renamed" {
		t.Errorf("expected renamed, got %s", results[0].Status)
	}
}

func TestRenameMany_KeyNotFound(t *testing.T) {
	entries := makeRenameMany()
	_, results := RenameMany(entries, map[string]string{"MISSING": "NEW"}, RenameManyOptions{})
	if results[0].Status != "not_found" {
		t.Errorf("expected not_found, got %s", results[0].Status)
	}
}

func TestRenameMany_ConflictWithoutOverwrite(t *testing.T) {
	entries := makeRenameMany()
	_, results := RenameMany(entries, map[string]string{"APP_HOST": "APP_PORT"}, RenameManyOptions{})
	if results[0].Status != "conflict" {
		t.Errorf("expected conflict, got %s", results[0].Status)
	}
}

func TestRenameMany_OverwritesOnConflict(t *testing.T) {
	entries := makeRenameMany()
	out, results := RenameMany(entries, map[string]string{"APP_HOST": "APP_PORT"}, RenameManyOptions{Overwrite: true})
	if results[0].Status != "renamed" {
		t.Errorf("expected renamed, got %s", results[0].Status)
	}
	v, _ := Lookup(out, "APP_PORT")
	if v != "localhost" {
		t.Errorf("expected localhost, got %s", v)
	}
}

func TestRenameMany_DryRun(t *testing.T) {
	entries := makeRenameMany()
	out, results := RenameMany(entries, map[string]string{"APP_HOST": "HOST"}, RenameManyOptions{DryRun: true})
	_, found := Lookup(out, "APP_HOST")
	if !found {
		t.Fatal("dry run should not modify entries")
	}
	if results[0].Status != "skipped" {
		t.Errorf("expected skipped, got %s", results[0].Status)
	}
}

func TestFormatRenameResults_ContainsArrow(t *testing.T) {
	results := []RenameResult{{OldKey: "A", NewKey: "B", Status: "renamed"}}
	out := FormatRenameResults(results, false)
	if !strings.Contains(out, "→") {
		t.Error("expected arrow in output")
	}
}

func TestFormatRenameSummaryLine_Counts(t *testing.T) {
	results := []RenameResult{
		{Status: "renamed"},
		{Status: "renamed"},
		{Status: "conflict"},
		{Status: "not_found"},
	}
	line := FormatRenameSummaryLine(results)
	if !strings.Contains(line, "renamed: 2") {
		t.Errorf("unexpected summary: %s", line)
	}
}
