package envfile

import (
	"testing"
)

func makeRenameEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
	}
}

func TestRename_BasicRename(t *testing.T) {
	entries := makeRenameEntries()
	out, res, err := Rename(entries, "APP_NAME", "APPLICATION_NAME", RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped {
		t.Fatalf("expected rename to succeed, got skipped: %s", res.Reason)
	}
	m := ToMap(out)
	if _, ok := m["APP_NAME"]; ok {
		t.Error("old key should not exist after rename")
	}
	if m["APPLICATION_NAME"] != "envoy" {
		t.Errorf("new key value mismatch: got %q", m["APPLICATION_NAME"])
	}
}

func TestRename_KeyNotFound(t *testing.T) {
	entries := makeRenameEntries()
	out, res, err := Rename(entries, "MISSING", "NEW_KEY", RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Skipped {
		t.Error("expected rename to be skipped when key is missing")
	}
	if len(out) != len(entries) {
		t.Error("entries should be unchanged")
	}
}

func TestRename_NewKeyExistsWithoutOverwrite(t *testing.T) {
	entries := makeRenameEntries()
	_, res, err := Rename(entries, "APP_NAME", "APP_ENV", RenameOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Skipped {
		t.Error("expected rename to be skipped when new key exists and Overwrite is false")
	}
}

func TestRename_NewKeyExistsWithOverwrite(t *testing.T) {
	entries := makeRenameEntries()
	out, res, err := Rename(entries, "APP_NAME", "APP_ENV", RenameOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped {
		t.Errorf("expected rename to succeed, got skipped: %s", res.Reason)
	}
	m := ToMap(out)
	if m["APP_ENV"] != "envoy" {
		t.Errorf("overwritten key should have old value: got %q", m["APP_ENV"])
	}
	if _, ok := m["APP_NAME"]; ok {
		t.Error("old key should not exist after rename")
	}
	// Length should shrink by one because APP_NAME was merged into APP_ENV.
	if len(out) != len(entries)-1 {
		t.Errorf("expected %d entries, got %d", len(entries)-1, len(out))
	}
}

func TestRename_DryRunDoesNotMutate(t *testing.T) {
	entries := makeRenameEntries()
	out, res, err := Rename(entries, "APP_NAME", "APPLICATION_NAME", RenameOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped {
		t.Error("dry-run result should not be marked skipped")
	}
	if ToMap(out)["APP_NAME"] != "envoy" {
		t.Error("dry-run should leave original entries unchanged")
	}
}

func TestRename_EmptyKeyReturnsError(t *testing.T) {
	entries := makeRenameEntries()
	_, _, err := Rename(entries, "", "NEW", RenameOptions{})
	if err == nil {
		t.Error("expected error for empty old key")
	}
}

func TestRename_SameKeySkipped(t *testing.T) {
	entries := makeRenameEntries()
	_, res, err := Rename(entries, "APP_NAME", "APP_NAME", RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Skipped {
		t.Error("expected rename to be skipped when old and new keys are identical")
	}
}
