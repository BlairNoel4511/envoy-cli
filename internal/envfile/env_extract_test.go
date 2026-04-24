package envfile

import "testing"

func makeExtractEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_PASSWORD", Value: "secret"},
		{Key: "DB_HOST", Value: "db.local"},
		{Key: "LOG_LEVEL", Value: "debug"},
	}
}

func TestExtract_ByPrefix(t *testing.T) {
	entries := makeExtractEntries()
	out, results, sum := Extract(entries, ExtractOptions{Prefix: "APP_"})
	if len(out) != 2 {
		t.Fatalf("expected 2 extracted, got %d", len(out))
	}
	if sum.Extracted != 2 {
		t.Errorf("expected summary.Extracted=2, got %d", sum.Extracted)
	}
	_ = results
}

func TestExtract_StripPrefix(t *testing.T) {
	entries := makeExtractEntries()
	out, _, _ := Extract(entries, ExtractOptions{Prefix: "APP_", StripPrefix: true})
	for _, e := range out {
		if e.Key == "APP_HOST" || e.Key == "APP_PORT" {
			t.Errorf("expected prefix stripped, got key %q", e.Key)
		}
	}
	keys := make(map[string]bool)
	for _, e := range out {
		keys[e.Key] = true
	}
	if !keys["HOST"] || !keys["PORT"] {
		t.Errorf("expected keys HOST and PORT after stripping, got %v", keys)
	}
}

func TestExtract_SensitiveOnly(t *testing.T) {
	entries := makeExtractEntries()
	out, _, sum := Extract(entries, ExtractOptions{SensitiveOnly: true})
	for _, e := range out {
		if !IsSensitive(e.Key) {
			t.Errorf("expected only sensitive keys, got %q", e.Key)
		}
	}
	if sum.Extracted == 0 {
		t.Error("expected at least one sensitive key extracted")
	}
}

func TestExtract_ByAllowlist(t *testing.T) {
	entries := makeExtractEntries()
	out, results, sum := Extract(entries, ExtractOptions{Keys: []string{"APP_HOST", "LOG_LEVEL"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 extracted, got %d", len(out))
	}
	if sum.Skipped != 3 {
		t.Errorf("expected 3 skipped, got %d", sum.Skipped)
	}
	_ = results
}

func TestExtract_EmptyOptions_ReturnsAll(t *testing.T) {
	entries := makeExtractEntries()
	out, _, sum := Extract(entries, ExtractOptions{})
	if len(out) != len(entries) {
		t.Errorf("expected all %d entries, got %d", len(entries), len(out))
	}
	if sum.Extracted != len(entries) {
		t.Errorf("expected summary.Extracted=%d, got %d", len(entries), sum.Extracted)
	}
}

func TestExtract_ResultContainsExtractedAs(t *testing.T) {
	entries := makeExtractEntries()
	_, results, _ := Extract(entries, ExtractOptions{Prefix: "DB_", StripPrefix: true})
	for _, r := range results {
		if r.Key == "DB_HOST" && r.ExtractedAs != "HOST" {
			t.Errorf("expected ExtractedAs=HOST, got %q", r.ExtractedAs)
		}
	}
}
