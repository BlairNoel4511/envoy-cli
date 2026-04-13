package envfile

import (
	"strings"
	"testing"
)

func makeCloneEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestClone_CopiesAllKeys(t *testing.T) {
	src := makeCloneEntries("APP_NAME", "envoy", "APP_PORT", "8080")
	dst := []Entry{}
	out, result := Clone(src, dst, CloneOptions{})
	if len(result.Cloned) != 2 {
		t.Fatalf("expected 2 cloned, got %d", len(result.Cloned))
	}
	m := ToMap(out)
	if m["APP_NAME"] != "envoy" {
		t.Errorf("expected APP_NAME=envoy")
	}
}

func TestClone_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := makeCloneEntries("APP_NAME", "new")
	dst := makeCloneEntries("APP_NAME", "old")
	_, result := Clone(src, dst, CloneOptions{Overwrite: false})
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
}

func TestClone_OverwritesExistingWhenEnabled(t *testing.T) {
	src := makeCloneEntries("APP_NAME", "new")
	dst := makeCloneEntries("APP_NAME", "old")
	out, result := Clone(src, dst, CloneOptions{Overwrite: true})
	if len(result.Cloned) != 1 {
		t.Errorf("expected 1 cloned, got %d", len(result.Cloned))
	}
	m := ToMap(out)
	if m["APP_NAME"] != "new" {
		t.Errorf("expected APP_NAME=new after overwrite, got %s", m["APP_NAME"])
	}
}

func TestClone_SkipsSensitiveKeys(t *testing.T) {
	src := makeCloneEntries("SECRET_KEY", "abc123", "APP_PORT", "9000")
	dst := []Entry{}
	_, result := Clone(src, dst, CloneOptions{SkipSensitive: true})
	for _, k := range result.Cloned {
		if IsSensitive(k) {
			t.Errorf("sensitive key %q should not have been cloned", k)
		}
	}
}

func TestClone_FiltersByPrefix(t *testing.T) {
	src := makeCloneEntries("APP_NAME", "envoy", "DB_HOST", "localhost")
	dst := []Entry{}
	out, result := Clone(src, dst, CloneOptions{Prefix: "APP_"})
	if len(result.Cloned) != 1 {
		t.Errorf("expected 1 cloned, got %d", len(result.Cloned))
	}
	m := ToMap(out)
	if _, ok := m["DB_HOST"]; ok {
		t.Errorf("DB_HOST should not be present")
	}
}

func TestFormatCloneResult_ContainsSummary(t *testing.T) {
	r := CloneResult{
		Cloned:  []string{"A", "B"},
		Skipped: []string{"C"},
	}
	s := FormatCloneResult(r)
	if !strings.Contains(s, "2") {
		t.Errorf("expected cloned count in output")
	}
	if !strings.Contains(s, "1") {
		t.Errorf("expected skipped count in output")
	}
}
