package envfile

import (
	"strings"
	"testing"
)

func makeScopeEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestScopeStore_SetAndGet(t *testing.T) {
	store := NewScopeStore()
	entries := makeScopeEntries("DB_HOST", "localhost", "DB_PORT", "5432")
	store.Set("dev", entries)

	got, ok := store.Get("dev")
	if !ok {
		t.Fatal("expected scope 'dev' to exist")
	}
	if len(got) != 2 {
		t.Errorf("expected 2 entries, got %d", len(got))
	}
}

func TestScopeStore_GetMissing(t *testing.T) {
	store := NewScopeStore()
	_, ok := store.Get("prod")
	if ok {
		t.Error("expected missing scope to return false")
	}
}

func TestScopeStore_Remove(t *testing.T) {
	store := NewScopeStore()
	store.Set("staging", makeScopeEntries("KEY", "val"))
	store.Remove("staging")
	_, ok := store.Get("staging")
	if ok {
		t.Error("expected scope to be removed")
	}
}

func TestScopeStore_List(t *testing.T) {
	store := NewScopeStore()
	store.Set("prod", makeScopeEntries("A", "1"))
	store.Set("dev", makeScopeEntries("B", "2"))
	store.Set("staging", makeScopeEntries("C", "3"))

	list := store.List()
	if len(list) != 3 {
		t.Fatalf("expected 3 scopes, got %d", len(list))
	}
	if list[0] != "dev" || list[1] != "prod" || list[2] != "staging" {
		t.Errorf("unexpected order: %v", list)
	}
}

func TestScopeStore_SetMeta(t *testing.T) {
	store := NewScopeStore()
	store.Set("dev", makeScopeEntries("X", "y"))

	err := store.SetMeta("dev", "owner", "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.Scopes["dev"].Meta["owner"] != "alice" {
		t.Error("expected meta to be set")
	}
}

func TestScopeStore_SetMeta_MissingScope(t *testing.T) {
	store := NewScopeStore()
	err := store.SetMeta("ghost", "k", "v")
	if err == nil {
		t.Error("expected error for missing scope")
	}
}

func TestFormatScopeList_Empty(t *testing.T) {
	store := NewScopeStore()
	out := FormatScopeList(store)
	if !strings.Contains(out, "No scopes") {
		t.Errorf("expected empty message, got: %q", out)
	}
}

func TestFormatScopeList_ShowsNames(t *testing.T) {
	store := NewScopeStore()
	store.Set("dev", makeScopeEntries("A", "1", "B", "2"))
	store.Set("prod", makeScopeEntries("C", "3"))

	out := FormatScopeList(store)
	if !strings.Contains(out, "dev") {
		t.Error("expected 'dev' in output")
	}
	if !strings.Contains(out, "prod") {
		t.Error("expected 'prod' in output")
	}
	if !strings.Contains(out, "2 key(s)") {
		t.Error("expected key count for dev")
	}
}
