package envfile

import (
	"strings"
	"testing"
)

func makeListEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "PORT", Value: "8080"},
		{Key: "API_SECRET", Value: "topsecret"},
		{Key: "DEBUG", Value: "true"},
	}
}

func TestList_ReturnsAllEntries(t *testing.T) {
	entries := makeListEntries()
	results := List(entries, ListOptions{})
	if len(results) != len(entries) {
		t.Fatalf("expected %d results, got %d", len(entries), len(results))
	}
}

func TestList_FilterByPrefix(t *testing.T) {
	results := List(makeListEntries(), ListOptions{FilterPrefix: "DB_"})
	if len(results) != 1 || results[0].Key != "DB_PASSWORD" {
		t.Fatalf("expected only DB_PASSWORD, got %+v", results)
	}
}

func TestList_FilterByPrefix_NoMatch(t *testing.T) {
	results := List(makeListEntries(), ListOptions{FilterPrefix: "NONEXISTENT_"})
	if len(results) != 0 {
		t.Fatalf("expected no results for unmatched prefix, got %+v", results)
	}
}

func TestList_RedactsSensitiveKeys(t *testing.T) {
	results := List(makeListEntries(), ListOptions{RedactSecrets: true})
	for _, r := range results {
		if IsSensitive(r.Key) {
			if r.Value != "***" {
				t.Errorf("expected %s to be redacted, got %q", r.Key, r.Value)
			}
			if !r.Masked {
				t.Errorf("expected Masked=true for %s", r.Key)
			}
		}
	}
}

func TestList_SortedOutput(t *testing.T) {
	results := List(makeListEntries(), ListOptions{SortKeys: true})
	for i := 1; i < len(results); i++ {
		if results[i].Key < results[i-1].Key {
			t.Errorf("results not sorted at index %d: %s < %s", i, results[i].Key, results[i-1].Key)
		}
	}
}

func TestList_IndexIsOneBasedSequential(t *testing.T) {
	results := List(makeListEntries(), ListOptions{})
	for i, r := range results {
		if r.Index != i+1 {
			t.Errorf("expected index %d, got %d", i+1, r.Index)
		}
	}
}

func TestListSummary_WithRedacted(t *testing.T) {
	results := List(makeListEntries(), ListOptions{RedactSecrets: true})
	summary := ListSummary(results)
	if !strings.Contains(summary, "redacted") {
		t.Errorf("expected summary to mention redacted, got %q", summary)
	}
}

func TestListSummary_NoRedacted(t *testing.T) {
	results := List(makeListEntries(), ListOptions{RedactSecrets: false})
	summary := ListSummary(results)
	if strings.Contains(summary, "redacted") {
		t.Errorf("did not expect redacted in summary, got %q", summary)
	}
}
