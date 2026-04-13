package envfile

import (
	"strings"
	"testing"
)

func TestDeprecationStore_DeprecateAndIsDeprecated(t *testing.T) {
	store := NewDeprecationStore()
	store.Deprecate("OLD_API_KEY", "use NEW_API_KEY", "NEW_API_KEY", DeprecationWarning)

	if !store.IsDeprecated("OLD_API_KEY") {
		t.Fatal("expected OLD_API_KEY to be deprecated")
	}
	if store.IsDeprecated("NEW_API_KEY") {
		t.Fatal("NEW_API_KEY should not be deprecated")
	}
}

func TestDeprecationStore_Get(t *testing.T) {
	store := NewDeprecationStore()
	store.Deprecate("LEGACY_TOKEN", "superseded", "", DeprecationRemoved)

	e, ok := store.Get("LEGACY_TOKEN")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Key != "LEGACY_TOKEN" {
		t.Errorf("unexpected key: %s", e.Key)
	}
	if e.Status != DeprecationRemoved {
		t.Errorf("unexpected status: %s", e.Status)
	}
}

func TestDeprecationStore_Remove(t *testing.T) {
	store := NewDeprecationStore()
	store.Deprecate("OLD_HOST", "no longer used", "", DeprecationActive)
	store.Remove("OLD_HOST")

	if store.IsDeprecated("OLD_HOST") {
		t.Fatal("expected OLD_HOST to be removed from store")
	}
}

func TestDeprecationStore_CheckEntries(t *testing.T) {
	store := NewDeprecationStore()
	store.Deprecate("OLD_DB_URL", "use DATABASE_URL", "DATABASE_URL", DeprecationWarning)

	entries := []Entry{
		{Key: "OLD_DB_URL", Value: "postgres://..."},
		{Key: "DATABASE_URL", Value: "postgres://..."},
	}
	hits := store.CheckEntries(entries)
	if len(hits) != 1 {
		t.Fatalf("expected 1 hit, got %d", len(hits))
	}
	if hits[0].Key != "OLD_DB_URL" {
		t.Errorf("unexpected hit key: %s", hits[0].Key)
	}
}

func TestFormatDeprecationList_Empty(t *testing.T) {
	store := NewDeprecationStore()
	out := FormatDeprecationList(store, false)
	if !strings.Contains(out, "no deprecated keys") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatDeprecationList_ContainsKey(t *testing.T) {
	store := NewDeprecationStore()
	store.Deprecate("OLD_SECRET", "replaced", "NEW_SECRET", DeprecationWarning)
	out := FormatDeprecationList(store, false)
	if !strings.Contains(out, "OLD_SECRET") {
		t.Errorf("expected OLD_SECRET in output, got: %s", out)
	}
	if !strings.Contains(out, "NEW_SECRET") {
		t.Errorf("expected replacement hint in output, got: %s", out)
	}
}

func TestFormatDeprecationList_Colorize(t *testing.T) {
	store := NewDeprecationStore()
	store.Deprecate("GONE_KEY", "removed entirely", "", DeprecationRemoved)
	out := FormatDeprecationList(store, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes in colorized output")
	}
}

func TestFormatDeprecationSummary(t *testing.T) {
	store := NewDeprecationStore()
	store.Deprecate("A", "old", "", DeprecationActive)
	store.Deprecate("B", "older", "", DeprecationRemoved)
	out := FormatDeprecationSummary(store)
	if !strings.Contains(out, "2") {
		t.Errorf("expected count 2 in summary, got: %s", out)
	}
}
