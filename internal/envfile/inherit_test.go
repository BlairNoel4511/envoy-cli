package envfile

import (
	"strings"
	"testing"
)

func makeInheritEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestInherit_AddsNewKeysFromParent(t *testing.T) {
	parent := makeInheritEntries("APP_ENV", "production", "LOG_LEVEL", "info")
	child := makeInheritEntries("APP_ENV", "development")

	out, result := Inherit(parent, child, InheritOptions{})
	if len(result.Added) != 1 || result.Added[0] != "LOG_LEVEL" {
		t.Errorf("expected LOG_LEVEL added, got %v", result.Added)
	}
	m := ToMap(out)
	if m["APP_ENV"] != "development" {
		t.Errorf("expected child value preserved, got %s", m["APP_ENV"])
	}
	if m["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %s", m["LOG_LEVEL"])
	}
}

func TestInherit_SkipsExistingWithoutOverwrite(t *testing.T) {
	parent := makeInheritEntries("APP_ENV", "production")
	child := makeInheritEntries("APP_ENV", "development")

	_, result := Inherit(parent, child, InheritOptions{Overwrite: false})
	if len(result.Skipped) != 1 || result.Skipped[0] != "APP_ENV" {
		t.Errorf("expected APP_ENV skipped, got %v", result.Skipped)
	}
}

func TestInherit_OverwritesWhenEnabled(t *testing.T) {
	parent := makeInheritEntries("APP_ENV", "production")
	child := makeInheritEntries("APP_ENV", "development")

	out, result := Inherit(parent, child, InheritOptions{Overwrite: true})
	if len(result.Overwritten) != 1 || result.Overwritten[0] != "APP_ENV" {
		t.Errorf("expected APP_ENV overwritten, got %v", result.Overwritten)
	}
	if ToMap(out)["APP_ENV"] != "production" {
		t.Error("expected overwritten value to be production")
	}
}

func TestInherit_SkipsSensitiveKeys(t *testing.T) {
	parent := makeInheritEntries("SECRET_KEY", "abc123", "APP_ENV", "prod")
	child := makeInheritEntries()

	_, result := Inherit(parent, child, InheritOptions{SkipSensitive: true})
	for _, k := range result.Skipped {
		if k == "SECRET_KEY" {
			return
		}
	}
	t.Error("expected SECRET_KEY to be skipped as sensitive")
}

func TestInherit_SkipsIdenticalValues(t *testing.T) {
	parent := makeInheritEntries("APP_ENV", "production")
	child := makeInheritEntries("APP_ENV", "production")

	_, result := Inherit(parent, child, InheritOptions{Overwrite: true})
	if len(result.Overwritten) != 0 {
		t.Errorf("expected no overwrites for identical values, got %v", result.Overwritten)
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected identical value to be skipped, got %v", result.Skipped)
	}
}

func TestFormatInheritResult_ContainsSummary(t *testing.T) {
	r := InheritResult{
		Added:       []string{"LOG_LEVEL"},
		Overwritten: []string{"APP_ENV"},
		Skipped:     []string{"SECRET_KEY"},
	}
	out := FormatInheritResult(r, false)
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Error("expected LOG_LEVEL in output")
	}
	if !strings.Contains(out, "Summary:") {
		t.Error("expected Summary in output")
	}
}
