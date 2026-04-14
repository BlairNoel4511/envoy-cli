package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAliasStore_AddAndGet(t *testing.T) {
	s := NewAliasStore()
	if err := s.Add("DB_URL", "DATABASE_URL", "legacy alias"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a, ok := s.Get("DB_URL")
	if !ok {
		t.Fatal("expected alias to exist")
	}
	if a.Canonical != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %s", a.Canonical)
	}
}

func TestAliasStore_AddSelfReference(t *testing.T) {
	s := NewAliasStore()
	if err := s.Add("KEY", "KEY", ""); err == nil {
		t.Fatal("expected error for self-referencing alias")
	}
}

func TestAliasStore_Remove(t *testing.T) {
	s := NewAliasStore()
	_ = s.Add("OLD", "NEW", "")
	if !s.Remove("OLD") {
		t.Fatal("expected Remove to return true")
	}
	_, ok := s.Get("OLD")
	if ok {
		t.Fatal("expected alias to be removed")
	}
	if s.Remove("OLD") {
		t.Fatal("expected Remove to return false for missing alias")
	}
}

func TestAliasStore_Resolve(t *testing.T) {
	s := NewAliasStore()
	_ = s.Add("DB_URL", "DATABASE_URL", "")
	if got := s.Resolve("DB_URL"); got != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %s", got)
	}
	if got := s.Resolve("UNKNOWN"); got != "UNKNOWN" {
		t.Errorf("expected UNKNOWN passthrough, got %s", got)
	}
}

func TestAliasStore_ResolveEntries(t *testing.T) {
	s := NewAliasStore()
	_ = s.Add("DB_URL", "DATABASE_URL", "")
	entries := []Entry{
		{Key: "DB_URL", Value: "postgres://localhost"},
		{Key: "PORT", Value: "5432"},
	}
	resolved := s.ResolveEntries(entries)
	if resolved[0].Key != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %s", resolved[0].Key)
	}
	if resolved[1].Key != "PORT" {
		t.Errorf("expected PORT unchanged, got %s", resolved[1].Key)
	}
}

func TestAliasStore_List_Sorted(t *testing.T) {
	s := NewAliasStore()
	_ = s.Add("Z_ALIAS", "Z_KEY", "")
	_ = s.Add("A_ALIAS", "A_KEY", "")
	list := s.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(list))
	}
	if list[0].Alias != "A_ALIAS" {
		t.Errorf("expected A_ALIAS first, got %s", list[0].Alias)
	}
}

func TestSaveAndLoadAliasStore_RoundTrip(t *testing.T) {
	s := NewAliasStore()
	_ = s.Add("DB_URL", "DATABASE_URL", "legacy")
	_ = s.Add("PG_PASS", "POSTGRES_PASSWORD", "")

	path := filepath.Join(t.TempDir(), "aliases.json")
	if err := SaveAliasStore(path, s); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	loaded, err := LoadAliasStore(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	a, ok := loaded.Get("DB_URL")
	if !ok || a.Canonical != "DATABASE_URL" {
		t.Errorf("round-trip failed for DB_URL")
	}
}

func TestLoadAliasStore_FileNotFound_ReturnsEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	s, err := LoadAliasStore(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.List()) != 0 {
		t.Error("expected empty store")
	}
}

func TestFormatAliasList_Empty(t *testing.T) {
	out := FormatAliasList([]Alias{})
	if !strings.Contains(out, "no aliases") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatAliasList_ContainsAlias(t *testing.T) {
	aliases := []Alias{{Alias: "DB_URL", Canonical: "DATABASE_URL", Comment: "legacy"}}
	out := FormatAliasList(aliases)
	if !strings.Contains(out, "DB_URL") || !strings.Contains(out, "DATABASE_URL") {
		t.Errorf("expected alias info in output, got: %s", out)
	}
	if !strings.Contains(out, "legacy") {
		t.Errorf("expected comment in output, got: %s", out)
	}
}

func TestSaveAliasStore_RestrictedPerms(t *testing.T) {
	s := NewAliasStore()
	_ = s.Add("X", "Y", "")
	path := filepath.Join(t.TempDir(), "aliases.json")
	if err := SaveAliasStore(path, s); err != nil {
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
