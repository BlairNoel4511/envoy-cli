package envfile

import (
	"testing"
)

func makeLockEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_PASS", Value: "oldpass"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestApplyLock_PinsNewKey(t *testing.T) {
	entries := makeLockEntries()
	lf := NewLockFile()
	lf.Pin("NEW_KEY", "newval", "ci", "")

	out, res, err := ApplyLock(entries, lf, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ToMap(out)
	if m["NEW_KEY"] != "newval" {
		t.Errorf("expected NEW_KEY=newval, got %q", m["NEW_KEY"])
	}
	if len(res.Pinned) != 1 || res.Pinned[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY in pinned list")
	}
}

func TestApplyLock_ConflictWithoutOverwrite(t *testing.T) {
	entries := makeLockEntries()
	lf := NewLockFile()
	lf.Pin("DB_PASS", "newpass", "ci", "")

	_, res, err := ApplyLock(entries, lf, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflict) != 1 || res.Conflict[0] != "DB_PASS" {
		t.Errorf("expected DB_PASS in conflict list")
	}
}

func TestApplyLock_OverwritesOnConflict(t *testing.T) {
	entries := makeLockEntries()
	lf := NewLockFile()
	lf.Pin("DB_PASS", "newpass", "ci", "")

	out, res, err := ApplyLock(entries, lf, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ToMap(out)
	if m["DB_PASS"] != "newpass" {
		t.Errorf("expected DB_PASS=newpass, got %q", m["DB_PASS"])
	}
	if len(res.Pinned) != 1 {
		t.Errorf("expected 1 pinned key")
	}
}

func TestApplyLock_SkipsIdenticalValues(t *testing.T) {
	entries := makeLockEntries()
	lf := NewLockFile()
	lf.Pin("PORT", "8080", "ci", "")

	_, res, err := ApplyLock(entries, lf, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "PORT" {
		t.Errorf("expected PORT in skipped list")
	}
}

func TestApplyLock_NilLockFileReturnsError(t *testing.T) {
	_, _, err := ApplyLock(makeLockEntries(), nil, false)
	if err == nil {
		t.Fatal("expected error for nil lock file")
	}
}

func TestValidateLock_NoViolations(t *testing.T) {
	entries := makeLockEntries()
	lf := NewLockFile()
	lf.Pin("PORT", "8080", "ci", "")

	v := ValidateLock(entries, lf)
	if len(v) != 0 {
		t.Errorf("expected no violations, got %v", v)
	}
}

func TestValidateLock_MissingKey(t *testing.T) {
	entries := makeLockEntries()
	lf := NewLockFile()
	lf.Pin("GHOST_KEY", "val", "ci", "")

	v := ValidateLock(entries, lf)
	if len(v) != 1 {
		t.Errorf("expected 1 violation, got %d", len(v))
	}
}

func TestValidateLock_WrongValue(t *testing.T) {
	entries := makeLockEntries()
	lf := NewLockFile()
	lf.Pin("APP_ENV", "staging", "ci", "")

	v := ValidateLock(entries, lf)
	if len(v) != 1 {
		t.Errorf("expected 1 violation, got %d", len(v))
	}
}
