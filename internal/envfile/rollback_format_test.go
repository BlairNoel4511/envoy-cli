package envfile

import (
	"strings"
	"testing"
)

func TestFormatRollbackPlan_NoEntries(t *testing.T) {
	plan := RollbackPlan{OperationID: "empty-op"}
	out := FormatRollbackPlan(plan, false)
	if !strings.Contains(out, "no changes to revert") {
		t.Errorf("expected 'no changes to revert', got: %s", out)
	}
}

func TestFormatRollbackPlan_ShowsRevertedKey(t *testing.T) {
	plan := RollbackPlan{
		OperationID: "op-show",
		Entries: []RollbackEntry{
			{Key: "DB_HOST", OldValue: "localhost", NewValue: "prod.host", HadKey: true},
		},
	}
	out := FormatRollbackPlan(plan, false)
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output")
	}
	if !strings.Contains(out, "->") {
		t.Errorf("expected arrow in revert output")
	}
}

func TestFormatRollbackPlan_ShowsRemovedKey(t *testing.T) {
	plan := RollbackPlan{
		OperationID: "op-remove",
		Entries: []RollbackEntry{
			{Key: "NEW_KEY", OldValue: "", NewValue: "some_val", HadKey: false},
		},
	}
	out := FormatRollbackPlan(plan, false)
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected NEW_KEY in output")
	}
	if !strings.Contains(out, "remove") {
		t.Errorf("expected 'remove' label for new key")
	}
}

func TestFormatRollbackPlan_ColorizeAddsEscapeCodes(t *testing.T) {
	plan := RollbackPlan{
		OperationID: "op-color",
		Entries: []RollbackEntry{
			{Key: "APP", OldValue: "v1", NewValue: "v2", HadKey: true},
		},
	}
	out := FormatRollbackPlan(plan, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes in colorized output")
	}
}

func TestFormatRollbackSummary_Counts(t *testing.T) {
	plan := RollbackPlan{
		OperationID: "op-summary",
		Entries: []RollbackEntry{
			{Key: "A", OldValue: "1", NewValue: "2", HadKey: true},
			{Key: "B", OldValue: "", NewValue: "new", HadKey: false},
			{Key: "C", OldValue: "x", NewValue: "y", HadKey: true},
		},
	}
	out := FormatRollbackSummary(plan)
	if !strings.Contains(out, "2 key(s) reverted") {
		t.Errorf("expected reverted count, got: %s", out)
	}
	if !strings.Contains(out, "1 key(s) removed") {
		t.Errorf("expected removed count, got: %s", out)
	}
}
