package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeFreezeEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_PASSWORD", Value: "secret"},
	}
}

func TestFreeze_MarksKeyAsFrozen(t *testing.T) {
	entries := makeFreezeEntries()
	store := NewFreezeStore()
	results, sum := Freeze(entries, store, []string{"APP_HOST"}, FreezeOptions{})

	if !store.IsFrozen("APP_HOST") {
		t.Fatal("expected APP_HOST to be frozen")
	}
	if len(results) != 1 || !results[0].Frozen {
		t.Fatal("expected one frozen result")
	}
	if sum.Frozen != 1 {
		t.Fatalf("expected Frozen=1, got %d", sum.Frozen)
	}
}

func TestFreeze_AlreadyFrozenWithoutForce(t *testing.T) {
	entries := makeFreezeEntries()
	store := NewFreezeStore()
	store.Freeze("APP_PORT", time.Now())

	results, sum := Freeze(entries, store, []string{"APP_PORT"}, FreezeOptions{Force: false})

	if !results[0].Already {
		t.Fatal("expected Already=true")
	}
	if sum.Already != 1 {
		t.Fatalf("expected Already=1, got %d", sum.Already)
	}
}

func TestFreeze_ForceRefreezesKey(t *testing.T) {
	entries := makeFreezeEntries()
	store := NewFreezeStore()
	old := time.Now().Add(-time.Hour)
	store.Freeze("APP_PORT", old)

	results, sum := Freeze(entries, store, []string{"APP_PORT"}, FreezeOptions{Force: true})

	if !results[0].Frozen {
		t.Fatal("expected Frozen=true after force")
	}
	at, _ := store.FrozenAt("APP_PORT")
	if !at.After(old) {
		t.Fatal("expected timestamp to be updated")
	}
	if sum.Frozen != 1 {
		t.Fatalf("expected Frozen=1, got %d", sum.Frozen)
	}
}

func TestFreeze_KeyNotFound(t *testing.T) {
	entries := makeFreezeEntries()
	store := NewFreezeStore()
	results, sum := Freeze(entries, store, []string{"MISSING_KEY"}, FreezeOptions{})

	if !results[0].NotFound {
		t.Fatal("expected NotFound=true")
	}
	if sum.NotFound != 1 {
		t.Fatalf("expected NotFound=1, got %d", sum.NotFound)
	}
}

func TestSaveAndLoadFreezeStore_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "freeze.json")

	store := NewFreezeStore()
	store.Freeze("APP_HOST", time.Now().UTC().Truncate(time.Second))

	if err := SaveFreezeStore(path, store); err != nil {
		t.Fatalf("SaveFreezeStore: %v", err)
	}

	loaded, err := LoadFreezeStore(path)
	if err != nil {
		t.Fatalf("LoadFreezeStore: %v", err)
	}
	if !loaded.IsFrozen("APP_HOST") {
		t.Fatal("expected APP_HOST to be frozen after round-trip")
	}
}

func TestLoadFreezeStore_FileNotFound_ReturnsEmpty(t *testing.T) {
	store, err := LoadFreezeStore("/nonexistent/freeze.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store == nil || len(store.Entries) != 0 {
		t.Fatal("expected empty store")
	}
}

func TestSaveFreezeStore_RestrictedPerms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "freeze.json")
	store := NewFreezeStore()

	if err := SaveFreezeStore(path, store); err != nil {
		t.Fatalf("SaveFreezeStore: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected 0600, got %v", info.Mode().Perm())
	}
}
