package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func makeRollbackEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
	}
}

func TestNewRollbackPlan_DetectsChanges(t *testing.T) {
	before := makeRollbackEntries()
	after := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "prod.db.example.com"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "NEW_KEY", Value: "new_value"},
	}

	plan := NewRollbackPlan("op-1", before, after)
	if len(plan.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(plan.Entries))
	}
}

func TestRollbackPlan_Apply_RevertsValues(t *testing.T) {
	before := makeRollbackEntries()
	after := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "prod.db.example.com"},
		{Key: "DB_PORT", Value: "5432"},
	}

	plan := NewRollbackPlan("op-2", before, after)
	result, err := plan.Apply(after)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := ToMap(result)
	if m["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost after rollback, got %q", m["DB_HOST"])
	}
}

func TestRollbackPlan_Apply_RemovesNewKeys(t *testing.T) {
	before := makeRollbackEntries()
	after := append(makeRollbackEntries(), Entry{Key: "EXTRA", Value: "extra_val"})

	plan := NewRollbackPlan("op-3", before, after)
	result, err := plan.Apply(after)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := ToMap(result)
	if _, ok := m["EXTRA"]; ok {
		t.Error("expected EXTRA to be removed after rollback")
	}
}

func TestSaveAndLoadRollbackPlan_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rollback.json")

	before := makeRollbackEntries()
	after := []Entry{{Key: "APP_NAME", Value: "changed"}, {Key: "DB_HOST", Value: "localhost"}, {Key: "DB_PORT", Value: "5432"}}
	plan := NewRollbackPlan("op-rt", before, after)

	if err := SaveRollbackPlan(path, plan); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := LoadRollbackPlan(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.OperationID != plan.OperationID {
		t.Errorf("operation ID mismatch: got %q", loaded.OperationID)
	}
	if len(loaded.Entries) != len(plan.Entries) {
		t.Errorf("entry count mismatch: got %d", len(loaded.Entries))
	}
}

func TestLoadRollbackPlan_FileNotFound_ReturnsEmpty(t *testing.T) {
	plan, err := LoadRollbackPlan("/nonexistent/rollback.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(plan.Entries) != 0 {
		t.Error("expected empty plan for missing file")
	}
}

func TestRollbackPlan_Apply_EmptyPlanReturnsError(t *testing.T) {
	plan := RollbackPlan{OperationID: "empty"}
	_, err := plan.Apply(makeRollbackEntries())
	if err == nil {
		t.Error("expected error for empty rollback plan")
	}
}

func TestSaveRollbackPlan_RestrictedPerms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rollback.json")
	plan := NewRollbackPlan("perm-test", makeRollbackEntries(), makeRollbackEntries())
	plan.Entries = append(plan.Entries, RollbackEntry{Key: "X", OldValue: "a", NewValue: "b", HadKey: true})

	if err := SaveRollbackPlan(path, plan); err != nil {
		t.Fatalf("save error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 perms, got %v", info.Mode().Perm())
	}
}

func TestFormatRollbackPlan_ContainsKeys(t *testing.T) {
	before := makeRollbackEntries()
	after := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "changed"},
		{Key: "DB_PORT", Value: "5432"},
	}
	plan := NewRollbackPlan("fmt-test", before, after)
	out := FormatRollbackPlan(plan, false)
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
}
