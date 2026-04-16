package envfile

import (
	"testing"
)

func makeGetEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestGet_FoundKey(t *testing.T) {
	entries := makeGetEntries()
	r := Get(entries, "APP_NAME", GetOptions{})
	if !r.Found || r.Value != "envoy" {
		t.Errorf("expected found=true value=envoy, got found=%v value=%s", r.Found, r.Value)
	}
}

func TestGet_MissingKeyReturnsDefault(t *testing.T) {
	entries := makeGetEntries()
	r := Get(entries, "MISSING", GetOptions{Default: "fallback"})
	if r.Found || r.Value != "fallback" {
		t.Errorf("expected found=false value=fallback, got found=%v value=%s", r.Found, r.Value)
	}
}

func TestGet_RedactsSensitiveKey(t *testing.T) {
	entries := makeGetEntries()
	r := Get(entries, "DB_PASSWORD", GetOptions{Redact: true})
	if r.Value != "***" || !r.Redacted {
		t.Errorf("expected redacted value, got %s", r.Value)
	}
}

func TestGet_NoRedactWhenDisabled(t *testing.T) {
	entries := makeGetEntries()
	r := Get(entries, "DB_PASSWORD", GetOptions{Redact: false})
	if r.Value != "s3cr3t" || r.Redacted {
		t.Errorf("expected plain value, got %s", r.Value)
	}
}

func TestGetMany_ReturnsAllResults(t *testing.T) {
	entries := makeGetEntries()
	results := GetMany(entries, []string{"APP_NAME", "PORT", "NOPE"}, GetOptions{})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[2].Found {
		t.Error("expected third result to be not found")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	results := []GetResult{
		{Found: true},
		{Found: true},
		{Found: false},
	}
	s := GetSummary(results)
	if s != "2 found, 1 missing" {
		t.Errorf("unexpected summary: %s", s)
	}
}
