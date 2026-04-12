package envfile

import (
	"testing"
)

func makeFilterEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_PASSWORD", Value: "secret123"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "SECRET_KEY", Value: "topsecret"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	entries := makeFilterEntries()
	result := FilterByPrefix(entries, "APP_")
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key != "APP_NAME" && e.Key != "APP_ENV" {
			t.Errorf("unexpected key: %s", e.Key)
		}
	}
}

func TestFilter_SensitiveOnly(t *testing.T) {
	entries := makeFilterEntries()
	result := FilterSensitive(entries)
	if len(result) == 0 {
		t.Fatal("expected at least one sensitive entry")
	}
	for _, e := range result {
		if !IsSensitive(e.Key) {
			t.Errorf("key %q should not be in sensitive results", e.Key)
		}
	}
}

func TestFilter_ByAllowlist(t *testing.T) {
	entries := makeFilterEntries()
	result := Filter(entries, FilterOptions{Keys: []string{"PORT", "DB_HOST"}})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestFilter_CombinedPrefixAndSensitive(t *testing.T) {
	entries := makeFilterEntries()
	result := Filter(entries, FilterOptions{Prefix: "DB_", SensitiveOnly: true})
	for _, e := range result {
		if !IsSensitive(e.Key) {
			t.Errorf("key %q should be sensitive", e.Key)
		}
		if e.Key[:3] != "DB_" {
			t.Errorf("key %q should have DB_ prefix", e.Key)
		}
	}
}

func TestFilter_EmptyPrefix_ReturnsAll(t *testing.T) {
	entries := makeFilterEntries()
	result := FilterByPrefix(entries, "")
	if len(result) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(result))
	}
}

func TestFilter_NoMatch_ReturnsNil(t *testing.T) {
	entries := makeFilterEntries()
	result := FilterByPrefix(entries, "NONEXISTENT_")
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}
