package envfile

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func makeProfileEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "development"},
		{Key: "DB_HOST", Value: "localhost"},
	}
}

func TestProfileStore_SetAndGet(t *testing.T) {
	ps := NewProfileStore()
	p := &Profile{Name: "dev", Entries: makeProfileEntries()}
	ps.Set(p)

	got, ok := ps.Get("dev")
	if !ok {
		t.Fatal("expected profile 'dev' to exist")
	}
	if got.Name != "dev" {
		t.Errorf("expected name 'dev', got %q", got.Name)
	}
	if len(got.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(got.Entries))
	}
}

func TestProfileStore_GetMissing(t *testing.T) {
	ps := NewProfileStore()
	_, ok := ps.Get("prod")
	if ok {
		t.Error("expected missing profile to return false")
	}
}

func TestProfileStore_Remove(t *testing.T) {
	ps := NewProfileStore()
	ps.Set(&Profile{Name: "staging", Entries: makeProfileEntries()})
	ps.Remove("staging")
	_, ok := ps.Get("staging")
	if ok {
		t.Error("expected profile to be removed")
	}
}

func TestProfileStore_List(t *testing.T) {
	ps := NewProfileStore()
	ps.Set(&Profile{Name: "dev"})
	ps.Set(&Profile{Name: "prod"})

	names := ps.List()
	if len(names) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(names))
	}
	if !slices.Contains(names, "dev") || !slices.Contains(names, "prod") {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestSaveAndLoadProfileStore_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profiles.json")

	ps := NewProfileStore()
	ps.Set(&Profile{Name: "dev", Entries: makeProfileEntries(), Tags: []string{"local"}})
	ps.Set(&Profile{Name: "prod", Entries: []Entry{{Key: "APP_ENV", Value: "production"}}})

	if err := SaveProfileStore(path, ps); err != nil {
		t.Fatalf("SaveProfileStore: %v", err)
	}

	loaded, err := LoadProfileStore(path)
	if err != nil {
		t.Fatalf("LoadProfileStore: %v", err)
	}
	if len(loaded.Profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(loaded.Profiles))
	}
	dev, ok := loaded.Get("dev")
	if !ok || len(dev.Tags) != 1 || dev.Tags[0] != "local" {
		t.Errorf("unexpected dev profile: %+v", dev)
	}
}

func TestLoadProfileStore_FileNotFound(t *testing.T) {
	ps, err := LoadProfileStore("/nonexistent/path/profiles.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(ps.Profiles) != 0 {
		t.Error("expected empty store for missing file")
	}
}

func TestSaveProfileStore_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profiles.json")
	ps := NewProfileStore()
	if err := SaveProfileStore(path, ps); err != nil {
		t.Fatalf("SaveProfileStore: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %v", info.Mode().Perm())
	}
}
