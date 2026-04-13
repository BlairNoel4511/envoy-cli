package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func makeGroupEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "secret"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "APP_ENV", Value: "production"},
	}
}

func TestGroupStore_AddAndGet(t *testing.T) {
	gs := NewGroupStore()
	gs.Add("database", []string{"DB_HOST", "DB_PASS"})
	g := gs.Get("database")
	if g == nil {
		t.Fatal("expected group to exist")
	}
	if len(g.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(g.Keys))
	}
}

func TestGroupStore_NoDuplicateKeys(t *testing.T) {
	gs := NewGroupStore()
	gs.Add("app", []string{"APP_PORT"})
	gs.Add("app", []string{"APP_PORT", "APP_ENV"})
	g := gs.Get("app")
	if len(g.Keys) != 2 {
		t.Errorf("expected 2 unique keys, got %d", len(g.Keys))
	}
}

func TestGroupStore_Remove(t *testing.T) {
	gs := NewGroupStore()
	gs.Add("tmp", []string{"X"})
	if err := gs.Remove("tmp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gs.Get("tmp") != nil {
		t.Error("expected group to be removed")
	}
}

func TestGroupStore_RemoveMissing(t *testing.T) {
	gs := NewGroupStore()
	if err := gs.Remove("ghost"); err == nil {
		t.Error("expected error for missing group")
	}
}

func TestGroupStore_List(t *testing.T) {
	gs := NewGroupStore()
	gs.Add("z-group", []string{"Z"})
	gs.Add("a-group", []string{"A"})
	names := gs.List()
	if names[0] != "a-group" || names[1] != "z-group" {
		t.Errorf("expected sorted names, got %v", names)
	}
}

func TestFilterByGroup_ReturnsMatchingEntries(t *testing.T) {
	gs := NewGroupStore()
	gs.Add("db", []string{"DB_HOST", "DB_PASS"})
	entries := makeGroupEntries()
	result, err := FilterByGroup(entries, gs, "db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestFilterByGroup_MissingGroup(t *testing.T) {
	gs := NewGroupStore()
	_, err := FilterByGroup(makeGroupEntries(), gs, "nope")
	if err == nil {
		t.Error("expected error for missing group")
	}
}

func TestSaveAndLoadGroupStore_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "groups.json")
	gs := NewGroupStore()
	gs.Add("infra", []string{"DB_HOST", "APP_PORT"})
	if err := SaveGroupStore(path, gs); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := LoadGroupStore(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	g := loaded.Get("infra")
	if g == nil || len(g.Keys) != 2 {
		t.Errorf("round-trip failed: %+v", loaded)
	}
}

func TestLoadGroupStore_FileNotFound_ReturnsEmpty(t *testing.T) {
	gs, err := LoadGroupStore(filepath.Join(os.TempDir(), "no_such_groups.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gs.Groups) != 0 {
		t.Error("expected empty store")
	}
}
