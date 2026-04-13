package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func makeArchiveEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestArchive_AddAndGet(t *testing.T) {
	a := NewArchive()
	record := a.Add("initial", makeArchiveEntries())

	got, ok := a.Get(record.ID)
	if !ok {
		t.Fatal("expected to find record by ID")
	}
	if got.Label != "initial" {
		t.Errorf("expected label 'initial', got %q", got.Label)
	}
	if len(got.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(got.Entries))
	}
}

func TestArchive_GetMissing(t *testing.T) {
	a := NewArchive()
	_, ok := a.Get("nonexistent-id")
	if ok {
		t.Error("expected not found for missing ID")
	}
}

func TestArchive_List(t *testing.T) {
	a := NewArchive()
	a.Add("v1", makeArchiveEntries())
	a.Add("v2", makeArchiveEntries())

	if len(a.List()) != 2 {
		t.Errorf("expected 2 records, got %d", len(a.List()))
	}
}

func TestSaveAndLoadArchive_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "archive.json")

	a := NewArchive()
	a.Add("snapshot-1", makeArchiveEntries())

	if err := SaveArchive(path, a); err != nil {
		t.Fatalf("SaveArchive failed: %v", err)
	}

	loaded, err := LoadArchive(path)
	if err != nil {
		t.Fatalf("LoadArchive failed: %v", err)
	}
	if len(loaded.Records) != 1 {
		t.Errorf("expected 1 record, got %d", len(loaded.Records))
	}
	if loaded.Records[0].Label != "snapshot-1" {
		t.Errorf("expected label 'snapshot-1', got %q", loaded.Records[0].Label)
	}
}

func TestLoadArchive_FileNotFound_ReturnsEmpty(t *testing.T) {
	a, err := LoadArchive("/nonexistent/path/archive.json")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(a.Records) != 0 {
		t.Errorf("expected empty archive, got %d records", len(a.Records))
	}
}

func TestSaveArchive_CreatesFileWithRestrictedPerms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "archive.json")

	a := NewArchive()
	a.Add("perms-test", makeArchiveEntries())

	if err := SaveArchive(path, a); err != nil {
		t.Fatalf("SaveArchive failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perms 0600, got %v", info.Mode().Perm())
	}
}
