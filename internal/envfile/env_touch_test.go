package envfile

import (
	"strings"
	"testing"
	"time"
)

func makeTouchEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_PASSWORD", Value: "secret", Sensitive: true},
	}
}

func TestTouch_SetsTimestampForAllKeys(t *testing.T) {
	entries := makeTouchEntries()
	store := NewTouchStore()
	at := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	results, summary := Touch(entries, store, TouchOptions{At: at})

	if summary.Touched != 3 {
		t.Fatalf("expected 3 touched, got %d", summary.Touched)
	}
	if summary.Skipped != 0 {
		t.Fatalf("expected 0 skipped, got %d", summary.Skipped)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.WasSet {
			t.Errorf("expected WasSet for key %s", r.Key)
		}
		ts, ok := store.Get(r.Key)
		if !ok || !ts.Equal(at) {
			t.Errorf("store timestamp mismatch for %s", r.Key)
		}
	}
}

func TestTouch_SkipsExistingWithoutOverwrite(t *testing.T) {
	entries := makeTouchEntries()
	store := NewTouchStore()
	first := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	store.Set("APP_HOST", first)

	_, summary := Touch(entries, store, TouchOptions{At: time.Now().UTC()})

	if summary.Skipped != 1 {
		t.Fatalf("expected 1 skipped, got %d", summary.Skipped)
	}
	ts, _ := store.Get("APP_HOST")
	if !ts.Equal(first) {
		t.Error("existing timestamp should not be overwritten")
	}
}

func TestTouch_OverwritesExistingWhenEnabled(t *testing.T) {
	entries := makeTouchEntries()
	store := NewTouchStore()
	old := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	newT := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	store.Set("APP_HOST", old)

	_, _ = Touch(entries, store, TouchOptions{At: newT, Overwrite: true})

	ts, ok := store.Get("APP_HOST")
	if !ok || !ts.Equal(newT) {
		t.Error("expected timestamp to be overwritten")
	}
}

func TestTouch_TargetSpecificKeys(t *testing.T) {
	entries := makeTouchEntries()
	store := NewTouchStore()
	at := time.Now().UTC()

	_, summary := Touch(entries, store, TouchOptions{At: at, Keys: []string{"APP_PORT"}})

	if summary.Touched != 1 {
		t.Fatalf("expected 1 touched, got %d", summary.Touched)
	}
	if _, ok := store.Get("APP_HOST"); ok {
		t.Error("APP_HOST should not have been touched")
	}
}

func TestFormatTouchResults_ShowsTouched(t *testing.T) {
	at := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	results := []TouchResult{
		{Key: "APP_HOST", NewTouch: at, WasSet: true},
	}
	out := FormatTouchResults(results, false)
	if !strings.Contains(out, "APP_HOST") {
		t.Error("expected APP_HOST in output")
	}
	if !strings.Contains(out, "2024-06-01") {
		t.Error("expected date in output")
	}
}

func TestFormatTouchResults_ColorizeAddsEscapeCodes(t *testing.T) {
	at := time.Now().UTC()
	results := []TouchResult{
		{Key: "X", NewTouch: at, WasSet: true},
	}
	out := FormatTouchResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI escape codes when colorize=true")
	}
}

func TestFormatTouchSummary_Counts(t *testing.T) {
	s := TouchSummary{Touched: 4, Skipped: 2}
	out := FormatTouchSummary(s)
	if !strings.Contains(out, "4") || !strings.Contains(out, "2") {
		t.Errorf("unexpected summary output: %s", out)
	}
}
