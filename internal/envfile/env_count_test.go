package envfile

import (
	"strings"
	"testing"
)

func makeCountEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "APP_SECRET", Value: "s3cr3t"},
		{Key: "DB_PASSWORD", Value: "hunter2"},
		{Key: "DEBUG", Value: ""},
		{Key: "LOG_LEVEL", Value: "info"},
	}
}

func TestCount_TotalEntries(t *testing.T) {
	entries := makeCountEntries()
	r := Count(entries, CountOptions{})
	if r.Total != 6 {
		t.Errorf("expected Total=6, got %d", r.Total)
	}
}

func TestCount_SensitiveEntries(t *testing.T) {
	entries := makeCountEntries()
	r := Count(entries, CountOptions{})
	// APP_SECRET and DB_PASSWORD should be sensitive
	if r.Sensitive < 1 {
		t.Errorf("expected at least 1 sensitive entry, got %d", r.Sensitive)
	}
}

func TestCount_EmptyValues(t *testing.T) {
	entries := makeCountEntries()
	r := Count(entries, CountOptions{})
	if r.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", r.Empty)
	}
	if r.NonEmpty != 5 {
		t.Errorf("expected NonEmpty=5, got %d", r.NonEmpty)
	}
}

func TestCount_WithPrefix(t *testing.T) {
	entries := makeCountEntries()
	r := Count(entries, CountOptions{Prefix: "APP_"})
	if r.Total != 3 {
		t.Errorf("expected Total=3 with prefix APP_, got %d", r.Total)
	}
}

func TestCount_FilteredBySensitive(t *testing.T) {
	entries := makeCountEntries()
	r := Count(entries, CountOptions{Sensitive: true})
	if r.Filtered != r.Sensitive {
		t.Errorf("expected Filtered=%d (sensitive), got %d", r.Sensitive, r.Filtered)
	}
}

func TestCount_FilteredByNonEmpty(t *testing.T) {
	entries := makeCountEntries()
	r := Count(entries, CountOptions{NonEmpty: true})
	if r.Filtered != r.NonEmpty {
		t.Errorf("expected Filtered=%d (non-empty), got %d", r.NonEmpty, r.Filtered)
	}
}

func TestCountSummary_ContainsTotal(t *testing.T) {
	entries := makeCountEntries()
	r := Count(entries, CountOptions{})
	summary := CountSummary(r)
	if !strings.Contains(summary, "total") {
		t.Errorf("expected summary to contain 'total', got: %s", summary)
	}
	if !strings.Contains(summary, "sensitive") {
		t.Errorf("expected summary to contain 'sensitive', got: %s", summary)
	}
}
