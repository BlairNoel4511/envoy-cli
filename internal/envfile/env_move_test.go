package envfile

import (
	"testing"
)

func makeMoveEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_URL", Value: "postgres://localhost/db"},
	}
}

func TestMove_BasicMove(t *testing.T) {
	entries := makeMoveEntries()
	updated, res, sum := Move(entries, "APP_HOST", "SERVICE_HOST", MoveOptions{})
	if res.Status != "moved" {
		t.Fatalf("expected moved, got %s", res.Status)
	}
	if sum.Moved != 1 {
		t.Fatalf("expected 1 moved")
	}
	m := ToMap(updated)
	if _, ok := m["APP_HOST"]; ok {
		t.Error("source key should be removed")
	}
	if m["SERVICE_HOST"] != "localhost" {
		t.Errorf("expected SERVICE_HOST=localhost, got %s", m["SERVICE_HOST"])
	}
}

func TestMove_SourceNotFound(t *testing.T) {
	entries := makeMoveEntries()
	_, res, sum := Move(entries, "MISSING", "NEW_KEY", MoveOptions{})
	if res.Status != "not_found" {
		t.Fatalf("expected not_found, got %s", res.Status)
	}
	if sum.NotFound != 1 {
		t.Fatalf("expected 1 not_found")
	}
}

func TestMove_ConflictWithoutOverwrite(t *testing.T) {
	entries := makeMoveEntries()
	_, res, sum := Move(entries, "APP_HOST", "APP_PORT", MoveOptions{Overwrite: false})
	if res.Status != "conflict" {
		t.Fatalf("expected conflict, got %s", res.Status)
	}
	if sum.Conflict != 1 {
		t.Fatalf("expected 1 conflict")
	}
}

func TestMove_OverwritesExistingWhenEnabled(t *testing.T) {
	entries := makeMoveEntries()
	updated, res, _ := Move(entries, "APP_HOST", "APP_PORT", MoveOptions{Overwrite: true})
	if res.Status != "moved" {
		t.Fatalf("expected moved, got %s", res.Status)
	}
	m := ToMap(updated)
	if m["APP_PORT"] != "localhost" {
		t.Errorf("expected APP_PORT=localhost after overwrite, got %s", m["APP_PORT"])
	}
	if _, ok := m["APP_HOST"]; ok {
		t.Error("source key should be removed")
	}
}

func TestMove_DryRunDoesNotModify(t *testing.T) {
	entries := makeMoveEntries()
	updated, res, _ := Move(entries, "APP_HOST", "SERVICE_HOST", MoveOptions{DryRun: true})
	if res.Status != "moved" {
		t.Fatalf("expected moved status in dry-run")
	}
	if res.Comment != "dry-run" {
		t.Fatalf("expected dry-run comment")
	}
	if len(updated) != len(entries) {
		t.Error("dry-run should not modify entries")
	}
}
