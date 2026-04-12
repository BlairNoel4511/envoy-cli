package envfile

import (
	"testing"
)

func TestMergeProfiles_AddsNewKeys(t *testing.T) {
	base := &Profile{Name: "dev", Entries: []Entry{
		{Key: "APP_ENV", Value: "development"},
	}}
	overlay := &Profile{Name: "overlay", Entries: []Entry{
		{Key: "DB_HOST", Value: "localhost"},
	}}

	result := MergeProfiles(base, overlay, ProfileMergeOptions{})
	if len(result.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result.Entries))
	}
}

func TestMergeProfiles_SkipsExistingWithoutOverwrite(t *testing.T) {
	base := &Profile{Name: "dev", Entries: []Entry{
		{Key: "APP_ENV", Value: "development"},
	}}
	overlay := &Profile{Name: "overlay", Entries: []Entry{
		{Key: "APP_ENV", Value: "CHANGED"},
	}}

	result := MergeProfiles(base, overlay, ProfileMergeOptions{Overwrite: false})
	m := ToMap(result.Entries)
	if m["APP_ENV"] != "development" {
		t.Errorf("expected original value, got %q", m["APP_ENV"])
	}
}

func TestMergeProfiles_OverwritesWhenEnabled(t *testing.T) {
	base := &Profile{Name: "dev", Entries: []Entry{
		{Key: "APP_ENV", Value: "development"},
	}}
	overlay := &Profile{Name: "overlay", Entries: []Entry{
		{Key: "APP_ENV", Value: "staging"},
	}}

	result := MergeProfiles(base, overlay, ProfileMergeOptions{Overwrite: true})
	m := ToMap(result.Entries)
	if m["APP_ENV"] != "staging" {
		t.Errorf("expected 'staging', got %q", m["APP_ENV"])
	}
}

func TestMergeProfiles_SkipSensitiveOnOverwrite(t *testing.T) {
	base := &Profile{Name: "dev", Entries: []Entry{
		{Key: "SECRET_KEY", Value: "original-secret"},
		{Key: "APP_ENV", Value: "development"},
	}}
	overlay := &Profile{Name: "overlay", Entries: []Entry{
		{Key: "SECRET_KEY", Value: "new-secret"},
		{Key: "APP_ENV", Value: "staging"},
	}}

	result := MergeProfiles(base, overlay, ProfileMergeOptions{Overwrite: true, SkipSensitive: true})
	m := ToMap(result.Entries)
	if m["SECRET_KEY"] != "original-secret" {
		t.Errorf("expected sensitive key to be preserved, got %q", m["SECRET_KEY"])
	}
	if m["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV to be overwritten, got %q", m["APP_ENV"])
	}
}

func TestMergeProfiles_PreservesBaseName(t *testing.T) {
	base := &Profile{Name: "production", Tags: []string{"live"}}
	overlay := &Profile{Name: "overlay"}
	result := MergeProfiles(base, overlay, ProfileMergeOptions{})
	if result.Name != "production" {
		t.Errorf("expected name 'production', got %q", result.Name)
	}
	if len(result.Tags) != 1 || result.Tags[0] != "live" {
		t.Errorf("expected tags to be preserved: %v", result.Tags)
	}
}

func TestProfileFromEntries_RoundTrip(t *testing.T) {
	entries := []Entry{{Key: "FOO", Value: "bar"}}
	p := ProfileFromEntries("test", entries)
	if p.Name != "test" || len(p.Entries) != 1 {
		t.Errorf("unexpected profile: %+v", p)
	}
}
