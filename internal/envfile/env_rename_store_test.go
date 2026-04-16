package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadRenameMapping_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rename.json")

	rm := RenameMapping{Mapping: map[string]string{
		"OLD_KEY": "NEW_KEY",
		"LEGACY":  "CURRENT",
	}}

	if err := SaveRenameMapping(path, rm); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := LoadRenameMapping(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.Mapping["OLD_KEY"] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY, got %s", loaded.Mapping["OLD_KEY"])
	}
	if loaded.Mapping["LEGACY"] != "CURRENT" {
		t.Errorf("expected CURRENT, got %s", loaded.Mapping["LEGACY"])
	}
}

func TestLoadRenameMapping_FileNotFound_ReturnsEmpty(t *testing.T) {
	loaded, err := LoadRenameMapping("/nonexistent/path/rename.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(loaded.Mapping) != 0 {
		t.Errorf("expected empty mapping")
	}
}

func TestSaveRenameMapping_RestrictedPerms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rename.json")
	rm := RenameMapping{Mapping: map[string]string{"A": "B"}}
	if err := SaveRenameMapping(path, rm); err != nil {
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
