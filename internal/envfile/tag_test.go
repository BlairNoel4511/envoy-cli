package envfile

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestTagStore_AddAndGet(t *testing.T) {
	ts := NewTagStore()
	ts.Add("DB_HOST", "infra")
	ts.Add("DB_HOST", "required")

	labels := ts.Get("DB_HOST")
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}
}

func TestTagStore_NoDuplicates(t *testing.T) {
	ts := NewTagStore()
	ts.Add("API_KEY", "secret")
	ts.Add("API_KEY", "secret")

	if len(ts.Get("API_KEY")) != 1 {
		t.Fatal("expected deduplication of identical labels")
	}
}

func TestTagStore_HasTag(t *testing.T) {
	ts := NewTagStore()
	ts.Add("PORT", "optional")

	if !ts.HasTag("PORT", "optional") {
		t.Error("expected HasTag to return true")
	}
	if ts.HasTag("PORT", "required") {
		t.Error("expected HasTag to return false for missing label")
	}
}

func TestTagStore_Remove(t *testing.T) {
	ts := NewTagStore()
	ts.Add("HOST", "infra")
	ts.Add("HOST", "required")
	ts.Remove("HOST", "infra")

	if ts.HasTag("HOST", "infra") {
		t.Error("expected label to be removed")
	}
	if !ts.HasTag("HOST", "required") {
		t.Error("expected remaining label to persist")
	}
}

func TestTagStore_KeysWithTag(t *testing.T) {
	ts := NewTagStore()
	ts.Add("DB_HOST", "infra")
	ts.Add("DB_PORT", "infra")
	ts.Add("API_KEY", "secret")

	keys := ts.KeysWithTag("infra")
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "DB_HOST" || keys[1] != "DB_PORT" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestSaveAndLoadTagStore_RoundTrip(t *testing.T) {
	ts := NewTagStore()
	ts.Add("DB_HOST", "infra")
	ts.Add("API_KEY", "secret")

	dir := t.TempDir()
	path := filepath.Join(dir, "tags.json")

	if err := SaveTagStore(path, ts); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := LoadTagStore(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if !loaded.HasTag("DB_HOST", "infra") {
		t.Error("expected DB_HOST infra tag after reload")
	}
	if !loaded.HasTag("API_KEY", "secret") {
		t.Error("expected API_KEY secret tag after reload")
	}
}

func TestLoadTagStore_FileNotFound_ReturnsEmpty(t *testing.T) {
	ts, err := LoadTagStore("/nonexistent/path/tags.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ts.All()) != 0 {
		t.Error("expected empty store for missing file")
	}
}

func TestSaveTagStore_RestrictedPerms(t *testing.T) {
	ts := NewTagStore()
	ts.Add("SECRET", "sensitive")

	dir := t.TempDir()
	path := filepath.Join(dir, "tags.json")
	if err := SaveTagStore(path, ts); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 perms, got %v", info.Mode().Perm())
	}
}
