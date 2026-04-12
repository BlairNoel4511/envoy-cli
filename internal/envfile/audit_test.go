package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAuditLog_RecordAndFilter(t *testing.T) {
	log := NewAuditLog()
	log.Record(AuditActionSet, "DB_PASS", "prod", true, "")
	log.Record(AuditActionSet, "APP_NAME", "prod", false, "")
	log.Record(AuditActionDelete, "OLD_KEY", "", false, "cleanup")

	if len(log.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(log.Entries))
	}

	sets := log.FilterByAction(AuditActionSet)
	if len(sets) != 2 {
		t.Errorf("expected 2 SET entries, got %d", len(sets))
	}

	deletes := log.FilterByAction(AuditActionDelete)
	if len(deletes) != 1 {
		t.Errorf("expected 1 DELETE entry, got %d", len(deletes))
	}
}

func TestAuditLog_FilterByKey(t *testing.T) {
	log := NewAuditLog()
	log.Record(AuditActionSet, "SECRET", "dev", true, "initial")
	log.Record(AuditActionSet, "OTHER", "dev", false, "")
	log.Record(AuditActionSync, "SECRET", "prod", true, "synced")

	results := log.FilterByKey("SECRET")
	if len(results) != 2 {
		t.Errorf("expected 2 entries for SECRET, got %d", len(results))
	}
}

func TestAuditEntrySummary_Format(t *testing.T) {
	log := NewAuditLog()
	log.Record(AuditActionExport, "API_KEY", "staging", true, "")

	summary := AuditEntrySummary(log.Entries[0])
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	for _, substr := range []string{"EXPORT", "API_KEY", "redacted", "staging"} {
		if !containsStr(summary, substr) {
			t.Errorf("expected summary to contain %q, got: %s", substr, summary)
		}
	}
}

func TestSaveAndLoadAuditLog_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")

	log := NewAuditLog()
	log.Record(AuditActionSet, "FOO", "default", false, "test")
	log.Record(AuditActionDelete, "BAR", "", false, "")

	if err := SaveAuditLog(path, log); err != nil {
		t.Fatalf("SaveAuditLog failed: %v", err)
	}

	loaded, err := LoadAuditLog(path)
	if err != nil {
		t.Fatalf("LoadAuditLog failed: %v", err)
	}
	if len(loaded.Entries) != 2 {
		t.Errorf("expected 2 entries after round-trip, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Key != "FOO" {
		t.Errorf("expected first key FOO, got %s", loaded.Entries[0].Key)
	}
}

func TestLoadAuditLog_FileNotFound_ReturnsEmpty(t *testing.T) {
	log, err := LoadAuditLog("/nonexistent/path/audit.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(log.Entries) != 0 {
		t.Errorf("expected empty log, got %d entries", len(log.Entries))
	}
}

func TestSaveAuditLog_RestrictedPerms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")

	log := NewAuditLog()
	log.Record(AuditActionSet, "KEY", "", false, "")

	if err := SaveAuditLog(path, log); err != nil {
		t.Fatalf("SaveAuditLog failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected perms 0600, got %v", info.Mode().Perm())
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
