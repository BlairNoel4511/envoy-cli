package envfile

import (
	"testing"
)

func makeUnsetEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_URL", Value: "postgres://localhost/db"},
	}
}

func TestUnset_RemovesExistingKey(t *testing.T) {
	entries := makeUnsetEntries()
	updated, res := Unset(entries, "APP_PORT")
	if !res.Removed {
		t.Fatal("expected Removed=true")
	}
	if res.Missing {
		t.Fatal("expected Missing=false")
	}
	for _, e := range updated {
		if e.Key == "APP_PORT" {
			t.Fatal("key still present after unset")
		}
	}
	if len(updated) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(updated))
	}
}

func TestUnset_MissingKey(t *testing.T) {
	entries := makeUnsetEntries()
	updated, res := Unset(entries, "NONEXISTENT")
	if res.Removed {
		t.Fatal("expected Removed=false")
	}
	if !res.Missing {
		t.Fatal("expected Missing=true")
	}
	if len(updated) != len(entries) {
		t.Fatal("entries should be unchanged")
	}
}

func TestUnsetMany_RemovesMultipleKeys(t *testing.T) {
	entries := makeUnsetEntries()
	updated, results, summary := UnsetMany(entries, []string{"APP_HOST", "DB_URL"})
	if summary.Removed != 2 {
		t.Fatalf("expected 2 removed, got %d", summary.Removed)
	}
	if summary.Missing != 0 {
		t.Fatalf("expected 0 missing, got %d", summary.Missing)
	}
	if len(updated) != 1 {
		t.Fatalf("expected 1 entry remaining, got %d", len(updated))
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestUnsetMany_SkipsMissingKeys(t *testing.T) {
	entries := makeUnsetEntries()
	_, _, summary := UnsetMany(entries, []string{"APP_HOST", "MISSING_KEY"})
	if summary.Removed != 1 {
		t.Fatalf("expected 1 removed, got %d", summary.Removed)
	}
	if summary.Missing != 1 {
		t.Fatalf("expected 1 missing, got %d", summary.Missing)
	}
}

func TestUnsetMany_EmptyKeys(t *testing.T) {
	entries := makeUnsetEntries()
	updated, results, summary := UnsetMany(entries, []string{})
	if len(updated) != len(entries) {
		t.Fatal("entries should be unchanged")
	}
	if len(results) != 0 {
		t.Fatal("expected no results")
	}
	if summary.Removed != 0 || summary.Missing != 0 {
		t.Fatal("expected zero counts")
	}
}
