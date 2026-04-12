package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLineage_RecordAndForKey(t *testing.T) {
	l := NewLineage()
	l.Record("set", "DB_HOST", "", "localhost", "manual")
	l.Record("set", "DB_PORT", "", "5432", "manual")
	l.Record("set", "DB_HOST", "localhost", "remotehost", "profile:prod")

	events := l.ForKey("DB_HOST")
	if len(events) != 2 {
		t.Fatalf("expected 2 events for DB_HOST, got %d", len(events))
	}
	if events[0].NewValue != "localhost" {
		t.Errorf("expected first new value 'localhost', got %q", events[0].NewValue)
	}
	if events[1].OldValue != "localhost" {
		t.Errorf("expected second old value 'localhost', got %q", events[1].OldValue)
	}
}

func TestLineageEvent_Summary(t *testing.T) {
	l := NewLineage()
	l.Record("set", "API_KEY", "", "abc123", "vault")
	l.Record("delete", "OLD_VAR", "val", "", "manual")
	l.Record("import", "NEW_VAR", "", "newval", "profile:staging")

	summaries := []string{}
	for _, e := range l.Events {
		summaries = append(summaries, e.Summary())
	}
	if len(summaries) != 3 {
		t.Fatalf("expected 3 summaries, got %d", len(summaries))
	}
	for _, s := range summaries {
		if s == "" {
			t.Error("expected non-empty summary")
		}
	}
}

func TestLineage_TrackDiff(t *testing.T) {
	l := NewLineage()
	d := DiffResult{
		Added:   []Entry{{Key: "NEW_KEY", Value: "newval"}},
		Removed: []Entry{{Key: "OLD_KEY", Value: "oldval"}},
		Changed: []ChangedEntry{{Key: "MOD_KEY", OldValue: "v1", NewValue: "v2"}},
	}
	l.TrackDiff(d, "profile:dev")

	if len(l.Events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(l.Events))
	}
	if l.Events[0].Action != "set" || l.Events[0].Key != "NEW_KEY" {
		t.Errorf("unexpected first event: %+v", l.Events[0])
	}
	if l.Events[1].Action != "delete" || l.Events[1].Key != "OLD_KEY" {
		t.Errorf("unexpected second event: %+v", l.Events[1])
	}
	if l.Events[2].Action != "set" || l.Events[2].Key != "MOD_KEY" {
		t.Errorf("unexpected third event: %+v", l.Events[2])
	}
}

func TestSaveAndLoadLineage_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lineage.json")

	l := NewLineage()
	l.Record("set", "FOO", "", "bar", "manual")

	if err := SaveLineage(path, l); err != nil {
		t.Fatalf("SaveLineage failed: %v", err)
	}

	loaded, err := LoadLineage(path)
	if err != nil {
		t.Fatalf("LoadLineage failed: %v", err)
	}
	if len(loaded.Events) != 1 || loaded.Events[0].Key != "FOO" {
		t.Errorf("unexpected loaded events: %+v", loaded.Events)
	}
}

func TestLoadLineage_FileNotFound_ReturnsEmpty(t *testing.T) {
	l, err := LoadLineage("/nonexistent/path/lineage.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(l.Events) != 0 {
		t.Errorf("expected empty lineage, got %d events", len(l.Events))
	}
}

func TestSaveLineage_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lineage.json")

	l := NewLineage()
	l.Record("set", "SECRET", "", "value", "vault")

	if err := SaveLineage(path, l); err != nil {
		t.Fatalf("SaveLineage failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected file perm 0600, got %v", info.Mode().Perm())
	}
}
