package envfile

import (
	"testing"
)

func TestDiff_Added(t *testing.T) {
	local := map[string]string{"FOO": "bar"}
	remote := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	entries := Diff(local, remote)
	var added []DiffEntry
	for _, e := range entries {
		if e.Type == DiffAdded {
			added = append(added, e)
		}
	}
	if len(added) != 1 || added[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY to be added, got %+v", added)
	}
}

func TestDiff_Removed(t *testing.T) {
	local := map[string]string{"FOO": "bar", "OLD_KEY": "old"}
	remote := map[string]string{"FOO": "bar"}

	entries := Diff(local, remote)
	var removed []DiffEntry
	for _, e := range entries {
		if e.Type == DiffRemoved {
			removed = append(removed, e)
		}
	}
	if len(removed) != 1 || removed[0].Key != "OLD_KEY" {
		t.Errorf("expected OLD_KEY to be removed, got %+v", removed)
	}
}

func TestDiff_Changed(t *testing.T) {
	local := map[string]string{"FOO": "old"}
	remote := map[string]string{"FOO": "new"}

	entries := Diff(local, remote)
	if len(entries) != 1 || entries[0].Type != DiffChanged {
		t.Errorf("expected FOO to be changed, got %+v", entries)
	}
	if entries[0].OldValue != "old" || entries[0].NewValue != "new" {
		t.Errorf("unexpected values: %+v", entries[0])
	}
}

func TestDiff_Unchanged(t *testing.T) {
	local := map[string]string{"FOO": "bar"}
	remote := map[string]string{"FOO": "bar"}

	entries := Diff(local, remote)
	if len(entries) != 1 || entries[0].Type != DiffUnchanged {
		t.Errorf("expected FOO to be unchanged, got %+v", entries)
	}
}

func TestHasChanges_True(t *testing.T) {
	entries := []DiffEntry{{Key: "FOO", Type: DiffAdded}}
	if !HasChanges(entries) {
		t.Error("expected HasChanges to return true")
	}
}

func TestHasChanges_False(t *testing.T) {
	entries := []DiffEntry{{Key: "FOO", Type: DiffUnchanged}}
	if HasChanges(entries) {
		t.Error("expected HasChanges to return false")
	}
}
