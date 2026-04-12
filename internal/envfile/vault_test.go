package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVault_AddAndGet(t *testing.T) {
	v := &Vault{}
	if err := v.Add("DB_PASS", "hunter2", "masterkey"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	val, err := v.Get("DB_PASS", "masterkey")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "hunter2" {
		t.Errorf("expected %q, got %q", "hunter2", val)
	}
}

func TestVault_GetMissingKey(t *testing.T) {
	v := &Vault{}
	_, err := v.Get("MISSING", "pass")
	if err != ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestVault_AddOverwritesExistingKey(t *testing.T) {
	v := &Vault{}
	v.Add("KEY", "old", "pass")
	v.Add("KEY", "new", "pass")
	if len(v.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(v.Entries))
	}
	val, _ := v.Get("KEY", "pass")
	if val != "new" {
		t.Errorf("expected %q, got %q", "new", val)
	}
}

func TestSaveAndLoadVault_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.json")

	v := &Vault{}
	v.Add("SECRET", "abc123", "passphrase")
	if err := SaveVault(path, v); err != nil {
		t.Fatalf("SaveVault failed: %v", err)
	}

	loaded, err := LoadVault(path)
	if err != nil {
		t.Fatalf("LoadVault failed: %v", err)
	}
	val, err := loaded.Get("SECRET", "passphrase")
	if err != nil {
		t.Fatalf("Get after load failed: %v", err)
	}
	if val != "abc123" {
		t.Errorf("expected %q, got %q", "abc123", val)
	}
}

func TestLoadVault_FileNotFound_ReturnsEmpty(t *testing.T) {
	v, err := LoadVault("/nonexistent/path/vault.json")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(v.Entries) != 0 {
		t.Errorf("expected empty vault, got %d entries", len(v.Entries))
	}
}

func TestSaveVault_RestrictedPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.json")
	v := &Vault{}
	v.Add("K", "v", "p")
	SaveVault(path, v)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %v", info.Mode().Perm())
	}
}
