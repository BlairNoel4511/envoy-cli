package envfile

import (
	"strings"
	"testing"
)

func makeResolveEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_URL", Value: "http://${APP_HOST}:${APP_PORT}/api"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "GREETING", Value: "hello $USER"},
		{Key: "PLAIN", Value: "no-vars-here"},
	}
}

func TestResolve_InterpolatesKnownVars(t *testing.T) {
	entries := makeResolveEntries()
	lookup := ToMap(entries)
	results := Resolve(entries, lookup, ResolveOption{})
	for _, r := range results {
		if r.Key == "APP_URL" {
			if !strings.Contains(r.Resolved, "localhost") {
				t.Errorf("expected APP_HOST resolved in APP_URL, got %q", r.Resolved)
			}
			if !strings.Contains(r.Resolved, "8080") {
				t.Errorf("expected APP_PORT resolved in APP_URL, got %q", r.Resolved)
			}
			if !r.Changed {
				t.Error("expected Changed=true for APP_URL")
			}
			return
		}
	}
	t.Error("APP_URL not found in results")
}

func TestResolve_TracksMissingVars(t *testing.T) {
	entries := makeResolveEntries()
	lookup := map[string]string{"APP_HOST": "localhost"} // APP_PORT missing
	results := Resolve(entries, lookup, ResolveOption{AllowMissing: true})
	for _, r := range results {
		if r.Key == "APP_URL" {
			if len(r.Unresolved) == 0 {
				t.Error("expected unresolved vars for APP_URL")
			}
			found := false
			for _, u := range r.Unresolved {
				if u == "APP_PORT" {
					found = true
				}
			}
			if !found {
				t.Errorf("expected APP_PORT in unresolved, got %v", r.Unresolved)
			}
			return
		}
	}
	t.Error("APP_URL not found in results")
}

func TestResolve_AllowMissingLeavesPlaceholder(t *testing.T) {
	entries := []Entry{{Key: "X", Value: "${MISSING_VAR}"}}
	results := Resolve(entries, map[string]string{}, ResolveOption{AllowMissing: true})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !strings.Contains(results[0].Resolved, "MISSING_VAR") {
		t.Errorf("expected placeholder preserved, got %q", results[0].Resolved)
	}
}

func TestResolve_DisallowMissingLeavesOriginal(t *testing.T) {
	original := "${UNKNOWN}"
	entries := []Entry{{Key: "X", Value: original}}
	results := Resolve(entries, map[string]string{}, ResolveOption{AllowMissing: false})
	if results[0].Resolved != original {
		t.Errorf("expected original value preserved, got %q", results[0].Resolved)
	}
	if results[0].Changed {
		t.Error("expected Changed=false when missing vars disallowed")
	}
}

func TestResolve_PrefixFilter(t *testing.T) {
	entries := makeResolveEntries()
	lookup := ToMap(entries)
	results := Resolve(entries, lookup, ResolveOption{Prefix: "APP_"})
	for _, r := range results {
		if !strings.HasPrefix(r.Key, "APP_") {
			t.Errorf("expected only APP_ keys, got %q", r.Key)
		}
	}
}

func TestApplyResolved_PatchesEntries(t *testing.T) {
	entries := makeResolveEntries()
	lookup := ToMap(entries)
	results := Resolve(entries, lookup, ResolveOption{})
	patched := ApplyResolved(entries, results)
	m := ToMap(patched)
	if !strings.Contains(m["APP_URL"], "localhost") {
		t.Errorf("expected patched APP_URL to contain localhost, got %q", m["APP_URL"])
	}
}

func TestResolve_PlainValueUnchanged(t *testing.T) {
	entries := []Entry{{Key: "PLAIN", Value: "no-vars-here"}}
	results := Resolve(entries, map[string]string{}, ResolveOption{})
	if results[0].Changed {
		t.Error("expected Changed=false for plain value")
	}
}
