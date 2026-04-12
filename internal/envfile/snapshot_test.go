package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSnapshot_ContainsAllKeys(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "localhost",
	}
	snap := NewSnapshot("test.env", env)

	if snap.Source != "test.env" {
		t.Errorf("expected source 'test.env', got %q", snap.Source)
	}
	if len(snap.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(snap.Entries))
	}
	if snap.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	env := map[string]string{
		"KEY_ONE": "value1",
		"KEY_TWO": "value2",
	}
	snap := NewSnapshot("roundtrip.env", env)

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "snap.json")

	if err := SaveSnapshot(path, snap); err != nil {
		t.Fatalf("SaveSnapshot failed: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot failed: %v", err)
	}

	if loaded.Source != snap.Source {
		t.Errorf("source mismatch: got %q, want %q", loaded.Source, snap.Source)
	}
	if loaded.Checksum != snap.Checksum {
		t.Errorf("checksum mismatch: got %q, want %q", loaded.Checksum, snap.Checksum)
	}
	if len(loaded.Entries) != len(snap.Entries) {
		t.Errorf("entry count mismatch: got %d, want %d", len(loaded.Entries), len(snap.Entries))
	}
}

func TestLoadSnapshot_FileNotFound(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSnapshot_ToMap(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	snap := NewSnapshot("test.env", env)
	result := snap.ToMap()

	for k, v := range env {
		if result[k] != v {
			t.Errorf("ToMap: key %q: got %q, want %q", k, result[k], v)
		}
	}
}

func TestSaveSnapshot_CreatesFileWithRestrictedPerms(t *testing.T) {
	snap := NewSnapshot("perms.env", map[string]string{"X": "1"})
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "snap.json")

	if err := SaveSnapshot(path, snap); err != nil {
		t.Fatalf("SaveSnapshot failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file mode 0600, got %v", info.Mode().Perm())
	}
}
