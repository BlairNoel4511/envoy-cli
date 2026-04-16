package envfile

import (
	"strings"
	"testing"
)

func makeSearchEntries() []Entry {
	return []Entry{
		{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
		{Key: "API_KEY", Value: "secret-abc"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DEBUG", Value: "false"},
	}
}

func TestSearch_MatchesKey(t *testing.T) {
	entries := makeSearchEntries()
	results := Search(entries, "API", SearchOptions{SearchKeys: true})
	if len(results) != 1 || results[0].Entry.Key != "API_KEY" {
		t.Errorf("expected API_KEY match, got %+v", results)
	}
}

func TestSearch_MatchesValue(t *testing.T) {
	entries := makeSearchEntries()
	results := Search(entries, "postgres", SearchOptions{SearchValues: true})
	if len(results) != 1 || results[0].Entry.Key != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL match")
	}
}

func TestSearch_CaseInsensitiveDefault(t *testing.T) {
	entries := makeSearchEntries()
	results := Search(entries, "production", SearchOptions{SearchValues: true, CaseSensitive: false})
	if len(results) != 1 {
		t.Errorf("expected 1 match, got %d", len(results))
	}
}

func TestSearch_RedactsSensitiveValues(t *testing.T) {
	entries := makeSearchEntries()
	results := Search(entries, "secret", SearchOptions{SearchValues: true, Redact: true})
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if results[0].Entry.Value != "***" {
		t.Errorf("expected redacted value, got %q", results[0].Entry.Value)
	}
}

func TestSearch_NoMatchReturnsEmpty(t *testing.T) {
	entries := makeSearchEntries()
	results := Search(entries, "NOTFOUND", SearchOptions{SearchKeys: true, SearchValues: true})
	if len(results) != 0 {
		t.Errorf("expected no results")
	}
}

func TestFormatSearchResults_Empty(t *testing.T) {
	out := FormatSearchResults(nil, false)
	if !strings.Contains(out, "no matches") {
		t.Errorf("expected no matches message")
	}
}

func TestFormatSearchSummary_Counts(t *testing.T) {
	s := SearchSummary{Query: "foo", Matched: 2, Total: 5}
	out := FormatSearchSummary(s)
	if !strings.Contains(out, "2/5") {
		t.Errorf("expected counts in summary, got %q", out)
	}
}
