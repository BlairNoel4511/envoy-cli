package envfile

import (
	"testing"
)

func makeDefaultEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "LOG_LEVEL", Value: ""},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestSetDefault_AppliesWhenEmpty(t *testing.T) {
	entries := makeDefaultEntries()
	result, res := SetDefault(entries, "LOG_LEVEL", "info", DefaultOptions{})
	if res.Status != "applied" {
		t.Fatalf("expected applied, got %s", res.Status)
	}
	m := ToMap(result)
	if m["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %s", m["LOG_LEVEL"])
	}
}

func TestSetDefault_SkipsExistingWithoutOverwrite(t *testing.T) {
	entries := makeDefaultEntries()
	_, res := SetDefault(entries, "APP_ENV", "staging", DefaultOptions{Overwrite: false})
	if res.Status != "skipped" {
		t.Fatalf("expected skipped, got %s", res.Status)
	}
}

func TestSetDefault_OverwritesExistingWhenEnabled(t *testing.T) {
	entries := makeDefaultEntries()
	result, res := SetDefault(entries, "APP_ENV", "staging", DefaultOptions{Overwrite: true})
	if res.Status != "applied" {
		t.Fatalf("expected applied, got %s", res.Status)
	}
	if ToMap(result)["APP_ENV"] != "staging" {
		t.Error("expected APP_ENV to be overwritten")
	}
}

func TestSetDefault_AppendsNewKey(t *testing.T) {
	entries := makeDefaultEntries()
	result, res := SetDefault(entries, "NEW_KEY", "hello", DefaultOptions{})
	if res.Status != "applied" {
		t.Fatalf("expected applied, got %s", res.Status)
	}
	if ToMap(result)["NEW_KEY"] != "hello" {
		t.Error("expected NEW_KEY to be appended")
	}
}

func TestSetDefault_UnchangedWhenValueIdentical(t *testing.T) {
	entries := makeDefaultEntries()
	_, res := SetDefault(entries, "APP_ENV", "production", DefaultOptions{Overwrite: true})
	if res.Status != "unchanged" {
		t.Fatalf("expected unchanged, got %s", res.Status)
	}
}

func TestSetDefault_SkipsSensitiveWithOption(t *testing.T) {
	entries := makeDefaultEntries()
	_, res := SetDefault(entries, "SECRET_KEY", "new", DefaultOptions{Overwrite: true, SkipSensitive: true})
	if res.Status != "skipped" {
		t.Fatalf("expected skipped for sensitive key, got %s", res.Status)
	}
}

func TestSetDefaults_BulkApply(t *testing.T) {
	entries := makeDefaultEntries()
	defaults := map[string]string{
		"LOG_LEVEL": "debug",
		"NEW_ONE":   "value1",
		"APP_ENV":   "staging",
	}
	_, _, sum := SetDefaults(entries, defaults, DefaultOptions{Overwrite: false})
	if sum.Applied < 2 {
		t.Errorf("expected at least 2 applied, got %d", sum.Applied)
	}
	if sum.Skipped < 1 {
		t.Errorf("expected at least 1 skipped, got %d", sum.Skipped)
	}
}
