package envfile

import (
	"testing"
)

func makeTrimEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "  myapp  "},
		{Key: "  DEBUG  ", Value: "true"},
		{Key: "DB_PASSWORD", Value: "  secret  "},
		{Key: "PORT", Value: "8080"},
	}
}

func TestTrim_TrimsValues(t *testing.T) {
	entries := makeTrimEntries()
	updated, results := Trim(entries, TrimOptions{TrimValues: true})

	if updated[0].Value != "myapp" {
		t.Errorf("expected 'myapp', got %q", updated[0].Value)
	}
	if !results[0].Changed {
		t.Error("expected result[0] to be marked changed")
	}
	if updated[3].Value != "8080" {
		t.Errorf("expected '8080' unchanged, got %q", updated[3].Value)
	}
	if results[3].Changed {
		t.Error("expected result[3] to be unchanged")
	}
}

func TestTrim_TrimsKeys(t *testing.T) {
	entries := makeTrimEntries()
	updated, results := Trim(entries, TrimOptions{TrimKeys: true})

	if updated[1].Key != "DEBUG" {
		t.Errorf("expected 'DEBUG', got %q", updated[1].Key)
	}
	if !results[1].Changed {
		t.Error("expected result[1] to be marked changed")
	}
}

func TestTrim_SkipsSensitiveValues(t *testing.T) {
	entries := makeTrimEntries()
	updated, results := Trim(entries, TrimOptions{
		TrimValues:    true,
		SkipSensitive: true,
	})

	// DB_PASSWORD is sensitive — should not be trimmed
	if updated[2].Value != "  secret  " {
		t.Errorf("expected sensitive value to be unchanged, got %q", updated[2].Value)
	}
	if !results[2].Skipped {
		t.Error("expected result[2] to be marked skipped")
	}
	if results[2].Changed {
		t.Error("expected result[2] not to be marked changed")
	}
}

func TestTrim_NoOpWhenBothFalse(t *testing.T) {
	entries := makeTrimEntries()
	updated, results := Trim(entries, TrimOptions{})

	for i, e := range updated {
		if e.Key != entries[i].Key || e.Value != entries[i].Value {
			t.Errorf("entry %d should be unchanged", i)
		}
		if results[i].Changed {
			t.Errorf("result %d should not be changed", i)
		}
	}
}

func TestTrim_BothKeysAndValues(t *testing.T) {
	entries := []Entry{
		{Key: "  HOST  ", Value: "  localhost  "},
	}
	updated, results := Trim(entries, TrimOptions{TrimKeys: true, TrimValues: true})

	if updated[0].Key != "HOST" {
		t.Errorf("expected key 'HOST', got %q", updated[0].Key)
	}
	if updated[0].Value != "localhost" {
		t.Errorf("expected value 'localhost', got %q", updated[0].Value)
	}
	if !results[0].Changed {
		t.Error("expected result to be changed")
	}
}

func TestTrimSummary_Counts(t *testing.T) {
	results := []TrimResult{
		{Changed: true},
		{Changed: true},
		{Skipped: true},
		{Changed: false},
	}
	summary := TrimSummary(results)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
