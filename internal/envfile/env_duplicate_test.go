package envfile

import (
	"testing"
)

func makeDuplicateEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
		{Key: "DEST_EXISTING", Value: "old"},
	}
}

func TestDuplicate_BasicCopy(t *testing.T) {
	entries := makeDuplicateEntries()
	result, _, sum := Duplicate(entries, "APP_NAME", "APP_NAME_COPY", DuplicateOptions{})
	if sum.Duplicated != 1 {
		t.Fatalf("expected 1 duplicated, got %d", sum.Duplicated)
	}
	v, ok := Lookup(result, "APP_NAME_COPY")
	if !ok || v != "envoy" {
		t.Errorf("expected APP_NAME_COPY=envoy, got %q", v)
	}
}

func TestDuplicate_SourceNotFound(t *testing.T) {
	entries := makeDuplicateEntries()
	_, res, sum := Duplicate(entries, "MISSING", "DEST", DuplicateOptions{})
	if res.Status != "source_not_found" {
		t.Errorf("expected source_not_found, got %q", res.Status)
	}
	if sum.NotFound != 1 {
		t.Errorf("expected NotFound=1, got %d", sum.NotFound)
	}
}

func TestDuplicate_SkipsExistingWithoutOverwrite(t *testing.T) {
	entries := makeDuplicateEntries()
	_, res, sum := Duplicate(entries, "APP_NAME", "DEST_EXISTING", DuplicateOptions{Overwrite: false})
	if res.Status != "skipped" {
		t.Errorf("expected skipped, got %q", res.Status)
	}
	if sum.Skipped != 1 {
		t.Errorf("expected Skipped=1, got %d", sum.Skipped)
	}
}

func TestDuplicate_OverwritesExistingWhenEnabled(t *testing.T) {
	entries := makeDuplicateEntries()
	result, res, _ := Duplicate(entries, "APP_NAME", "DEST_EXISTING", DuplicateOptions{Overwrite: true})
	if res.Status != "duplicated" {
		t.Errorf("expected duplicated, got %q", res.Status)
	}
	v, _ := Lookup(result, "DEST_EXISTING")
	if v != "envoy" {
		t.Errorf("expected DEST_EXISTING=envoy, got %q", v)
	}
}

func TestDuplicate_UnchangedWhenValueIdentical(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "same"},
		{Key: "B", Value: "same"},
	}
	_, res, sum := Duplicate(entries, "A", "B", DuplicateOptions{Overwrite: true})
	if res.Status != "unchanged" {
		t.Errorf("expected unchanged, got %q", res.Status)
	}
	if sum.Unchanged != 1 {
		t.Errorf("expected Unchanged=1, got %d", sum.Unchanged)
	}
}

func TestDuplicate_SkipsSensitiveKey(t *testing.T) {
	entries := makeDuplicateEntries()
	_, res, sum := Duplicate(entries, "SECRET_KEY", "SECRET_KEY_COPY", DuplicateOptions{SkipSensitive: true})
	if res.Status != "skipped" {
		t.Errorf("expected skipped for sensitive key, got %q", res.Status)
	}
	if sum.Skipped != 1 {
		t.Errorf("expected Skipped=1, got %d", sum.Skipped)
	}
}
