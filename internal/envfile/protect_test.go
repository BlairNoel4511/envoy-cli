package envfile

import (
	"strings"
	"testing"
)

func makeProtectEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "secret"},
		{Key: "APP_ENV", Value: "production"},
	}
}

func TestProtect_MarksKeyAsProtected(t *testing.T) {
	entries := makeProtectEntries()
	store := NewProtectStore()
	results, summary := Protect(entries, []string{"DB_HOST"}, store, ProtectOptions{})
	if summary.Protected != 1 {
		t.Fatalf("expected 1 protected, got %d", summary.Protected)
	}
	if !results[0].Protected {
		t.Error("expected Protected=true")
	}
	if !store.IsProtected("DB_HOST") {
		t.Error("expected DB_HOST to be in store")
	}
}

func TestProtect_AlreadyProtected_NoOverwrite(t *testing.T) {
	entries := makeProtectEntries()
	store := NewProtectStore()
	store.Add("DB_HOST")
	_, summary := Protect(entries, []string{"DB_HOST"}, store, ProtectOptions{Overwrite: false})
	if summary.Already != 1 {
		t.Fatalf("expected 1 already, got %d", summary.Already)
	}
}

func TestProtect_AlreadyProtected_WithOverwrite(t *testing.T) {
	entries := makeProtectEntries()
	store := NewProtectStore()
	store.Add("DB_HOST")
	_, summary := Protect(entries, []string{"DB_HOST"}, store, ProtectOptions{Overwrite: true})
	if summary.Protected != 1 {
		t.Fatalf("expected 1 protected on overwrite, got %d", summary.Protected)
	}
}

func TestProtect_KeyNotFound(t *testing.T) {
	entries := makeProtectEntries()
	store := NewProtectStore()
	_, summary := Protect(entries, []string{"MISSING_KEY"}, store, ProtectOptions{})
	if summary.NotFound != 1 {
		t.Fatalf("expected 1 not found, got %d", summary.NotFound)
	}
}

func TestSaveAndLoadProtectStore_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/protect.json"
	store := NewProtectStore()
	store.Add("DB_HOST")
	store.Add("APP_ENV")
	if err := SaveProtectStore(path, store); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	loaded, err := LoadProtectStore(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if !loaded.IsProtected("DB_HOST") || !loaded.IsProtected("APP_ENV") {
		t.Error("expected both keys to be protected after round-trip")
	}
}

func TestLoadProtectStore_FileNotFound_ReturnsEmpty(t *testing.T) {
	store, err := LoadProtectStore("/nonexistent/protect.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.All()) != 0 {
		t.Error("expected empty store")
	}
}

func TestFormatProtectResults_ShowsProtected(t *testing.T) {
	results := []ProtectResult{{Key: "DB_HOST", Protected: true}}
	out := FormatProtectResults(results, false)
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "protected") {
		t.Error("expected 'protected' label in output")
	}
}

func TestFormatProtectedList_SortsKeys(t *testing.T) {
	store := NewProtectStore()
	store.Add("Z_KEY")
	store.Add("A_KEY")
	out := FormatProtectedList(store)
	if strings.Index(out, "A_KEY") > strings.Index(out, "Z_KEY") {
		t.Error("expected A_KEY before Z_KEY in sorted output")
	}
}
