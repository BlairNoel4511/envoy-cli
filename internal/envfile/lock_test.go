package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLockFile_PinAndIsPinned(t *testing.T) {
	lf := NewLockFile()
	lf.Pin("API_KEY", "secret123", "alice", "pinned for prod")

	if !lf.IsPinned("API_KEY") {
		t.Fatal("expected API_KEY to be pinned")
	}
	if lf.IsPinned("MISSING") {
		t.Fatal("expected MISSING to not be pinned")
	}
}

func TestLockFile_Get(t *testing.T) {
	lf := NewLockFile()
	lf.Pin("DB_URL", "postgres://localhost/db", "bob", "")

	e, err := lf.Get("DB_URL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Value != "postgres://localhost/db" {
		t.Errorf("expected value %q, got %q", "postgres://localhost/db", e.Value)
	}
}

func TestLockFile_GetMissing(t *testing.T) {
	lf := NewLockFile()
	_, err := lf.Get("NOPE")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestLockFile_Unpin(t *testing.T) {
	lf := NewLockFile()
	lf.Pin("FOO", "bar", "alice", "")

	if !lf.Unpin("FOO") {
		t.Fatal("expected Unpin to return true")
	}
	if lf.IsPinned("FOO") {
		t.Fatal("expected FOO to be unpinned")
	}
	if lf.Unpin("FOO") {
		t.Fatal("expected Unpin to return false for already-removed key")
	}
}

func TestSaveAndLoadLockFile_RoundTrip(t *testing.T) {
	lf := NewLockFile()
	lf.Pin("SECRET", "abc", "ci", "automated")

	tmp := filepath.Join(t.TempDir(), "envoy.lock")
	if err := SaveLockFile(tmp, lf); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := LoadLockFile(tmp)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !loaded.IsPinned("SECRET") {
		t.Error("expected SECRET to be pinned after round-trip")
	}
}

func TestLoadLockFile_FileNotFound_ReturnsEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.lock")
	lf, err := LoadLockFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Entries) != 0 {
		t.Error("expected empty lock file")
	}
}

func TestSaveLockFile_CreatesFile(t *testing.T) {
	lf := NewLockFile()
	tmp := filepath.Join(t.TempDir(), "test.lock")
	if err := SaveLockFile(tmp, lf); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, err := os.Stat(tmp); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
