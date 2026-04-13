package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPinStore_AddAndIsPinned(t *testing.T) {
	store := NewPinStore()
	store.Add("DB_PASSWORD", "production", "sensitive")
	if !store.IsPinned("DB_PASSWORD", "production") {
		t.Error("expected DB_PASSWORD to be pinned to production")
	}
	if store.IsPinned("DB_PASSWORD", "staging") {
		t.Error("should not be pinned to staging")
	}
}

func TestPinStore_AddOverwritesReason(t *testing.T) {
	store := NewPinStore()
	store.Add("API_KEY", "production", "old reason")
	store.Add("API_KEY", "production", "new reason")
	pins := store.ForEnv("production")
	if len(pins) != 1 {
		t.Fatalf("expected 1 pin, got %d", len(pins))
	}
	if pins[0].Reason != "new reason" {
		t.Errorf("expected reason 'new reason', got %q", pins[0].Reason)
	}
}

func TestPinStore_Remove(t *testing.T) {
	store := NewPinStore()
	store.Add("SECRET", "prod", "")
	removed := store.Remove("SECRET", "prod")
	if !removed {
		t.Error("expected Remove to return true")
	}
	if store.IsPinned("SECRET", "prod") {
		t.Error("expected SECRET to be unpinned")
	}
}

func TestPinStore_RemoveMissing(t *testing.T) {
	store := NewPinStore()
	removed := store.Remove("MISSING", "prod")
	if removed {
		t.Error("expected Remove to return false for missing key")
	}
}

func TestPinStore_ForEnv(t *testing.T) {
	store := NewPinStore()
	store.Add("A", "prod", "")
	store.Add("B", "prod", "")
	store.Add("C", "staging", "")
	pins := store.ForEnv("prod")
	if len(pins) != 2 {
		t.Errorf("expected 2 pins for prod, got %d", len(pins))
	}
}

func TestSaveAndLoadPinStore_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	store := NewPinStore()
	store.Add("TOKEN", "production", "do not promote")
	if err := SavePinStore(path, store); err != nil {
		t.Fatalf("SavePinStore error: %v", err)
	}
	loaded, err := LoadPinStore(path)
	if err != nil {
		t.Fatalf("LoadPinStore error: %v", err)
	}
	if !loaded.IsPinned("TOKEN", "production") {
		t.Error("expected TOKEN to be pinned after round-trip")
	}
}

func TestLoadPinStore_FileNotFound_ReturnsEmpty(t *testing.T) {
	store, err := LoadPinStore("/nonexistent/path/pins.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Pins) != 0 {
		t.Error("expected empty pin store")
	}
}

func TestSavePinStore_CreatesFileWithRestrictedPerms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	store := NewPinStore()
	if err := SavePinStore(path, store); err != nil {
		t.Fatalf("SavePinStore error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perms 0600, got %v", info.Mode().Perm())
	}
}
