package envfile

import (
	"strings"
	"testing"
)

func TestFormatPatchResults_ShowsApplied(t *testing.T) {
	results := []PatchResult{
		{
			Instruction: PatchInstruction{Op: PatchSet, Key: "FOO", Value: "bar"},
			Applied:     true,
		},
	}
	out := FormatPatchResults(results, false)
	if !strings.Contains(out, "applied") {
		t.Errorf("expected 'applied' in output, got: %s", out)
	}
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected key in output, got: %s", out)
	}
}

func TestFormatPatchResults_ShowsSkipped(t *testing.T) {
	results := []PatchResult{
		{
			Instruction: PatchInstruction{Op: PatchDelete, Key: "MISSING"},
			Skipped:     true,
			Reason:      "key not found",
		},
	}
	out := FormatPatchResults(results, false)
	if !strings.Contains(out, "skipped") {
		t.Errorf("expected 'skipped' in output, got: %s", out)
	}
	if !strings.Contains(out, "key not found") {
		t.Errorf("expected reason in output, got: %s", out)
	}
}

func TestFormatPatchResults_ColorizeAddsEscapeCodes(t *testing.T) {
	results := []PatchResult{
		{
			Instruction: PatchInstruction{Op: PatchSet, Key: "X", Value: "1"},
			Applied:     true,
		},
	}
	out := FormatPatchResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes in colorized output")
	}
}

func TestFormatPatchSummary_Counts(t *testing.T) {
	results := []PatchResult{
		{Applied: true},
		{Applied: true},
		{Skipped: true, Reason: "identical value"},
	}
	summary := FormatPatchSummary(results)
	if !strings.Contains(summary, "2 applied") {
		t.Errorf("expected '2 applied', got: %s", summary)
	}
	if !strings.Contains(summary, "1 skipped") {
		t.Errorf("expected '1 skipped', got: %s", summary)
	}
}

func TestFormatPatchResults_RenameShowsArrow(t *testing.T) {
	results := []PatchResult{
		{
			Instruction: PatchInstruction{Op: PatchRename, Key: "OLD", NewKey: "NEW"},
			Applied:     true,
		},
	}
	out := FormatPatchResults(results, false)
	if !strings.Contains(out, "->") {
		t.Errorf("expected '->' in rename output, got: %s", out)
	}
	if !strings.Contains(out, "NEW") {
		t.Errorf("expected new key name in output, got: %s", out)
	}
}
