package envfile

import (
	"strings"
	"testing"
)

func makePromoteEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestPromote_AddsNewKeys(t *testing.T) {
	src := makePromoteEntries("APP_ENV", "production", "LOG_LEVEL", "info")
	dst := makePromoteEntries("EXISTING", "yes")

	result, summary := Promote(src, dst, PromoteOptions{})
	m := ToMap(result)

	if m["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", m["APP_ENV"])
	}
	if m["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", m["LOG_LEVEL"])
	}
	if !summary.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestPromote_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := makePromoteEntries("APP_ENV", "production")
	dst := makePromoteEntries("APP_ENV", "staging")

	result, summary := Promote(src, dst, PromoteOptions{Overwrite: false})
	m := ToMap(result)

	if m["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV to remain staging, got %q", m["APP_ENV"])
	}
	if summary.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestPromote_OverwritesWhenEnabled(t *testing.T) {
	src := makePromoteEntries("APP_ENV", "production")
	dst := makePromoteEntries("APP_ENV", "staging")

	result, summary := Promote(src, dst, PromoteOptions{Overwrite: true})
	m := ToMap(result)

	if m["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", m["APP_ENV"])
	}
	if !summary.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestPromote_SkipsSensitiveKeys(t *testing.T) {
	src := makePromoteEntries("SECRET_KEY", "abc123", "APP_ENV", "prod")
	dst := makePromoteEntries()

	result, summary := Promote(src, dst, PromoteOptions{SkipSensitive: true})
	m := ToMap(result)

	if _, ok := m["SECRET_KEY"]; ok {
		t.Error("expected SECRET_KEY to be skipped")
	}
	if m["APP_ENV"] != "prod" {
		t.Errorf("expected APP_ENV=prod, got %q", m["APP_ENV"])
	}
	_ = summary
}

func TestPromote_DryRunDoesNotWrite(t *testing.T) {
	src := makePromoteEntries("NEW_KEY", "value")
	dst := makePromoteEntries()

	result, summary := Promote(src, dst, PromoteOptions{DryRun: true})
	m := ToMap(result)

	if _, ok := m["NEW_KEY"]; ok {
		t.Error("expected NEW_KEY not to be written in dry run")
	}
	if !summary.DryRun {
		t.Error("expected DryRun flag to be set")
	}
}

func TestFormatPromoteSummary_ContainsCounts(t *testing.T) {
	src := makePromoteEntries("A", "1", "B", "2", "C", "3")
	dst := makePromoteEntries("B", "old")

	_, summary := Promote(src, dst, PromoteOptions{Overwrite: true})
	out := FormatPromoteSummary(summary)

	if !strings.Contains(out, "added") {
		t.Errorf("expected 'added' in output, got: %s", out)
	}
	if !strings.Contains(out, "overwritten") {
		t.Errorf("expected 'overwritten' in output, got: %s", out)
	}
}
