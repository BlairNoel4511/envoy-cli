package envfile

import (
	"strings"
	"testing"
)

func TestFormatResolveResults_ShowsChanged(t *testing.T) {
	results := []ResolveResult{
		{Key: "URL", Original: "http://${HOST}", Resolved: "http://localhost", Changed: true},
	}
	out := FormatResolveResults(results, false)
	if !strings.Contains(out, "URL") {
		t.Error("expected key URL in output")
	}
	if !strings.Contains(out, "http://localhost") {
		t.Error("expected resolved value in output")
	}
}

func TestFormatResolveResults_ShowsUnresolved(t *testing.T) {
	results := []ResolveResult{
		{Key: "X", Original: "${MISSING}", Resolved: "${MISSING}", Changed: false, Unresolved: []string{"MISSING"}},
	}
	out := FormatResolveResults(results, false)
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected unresolved var in output, got: %q", out)
	}
}

func TestFormatResolveResults_ColorizeAddsEscapeCodes(t *testing.T) {
	results := []ResolveResult{
		{Key: "A", Original: "$X", Resolved: "hello", Changed: true},
	}
	out := FormatResolveResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI escape codes in colorized output")
	}
}

func TestFormatResolveResults_EmptyWhenNoChanges(t *testing.T) {
	results := []ResolveResult{
		{Key: "A", Original: "plain", Resolved: "plain", Changed: false},
	}
	out := FormatResolveResults(results, false)
	if !strings.Contains(out, "all variables resolved") {
		t.Errorf("expected 'all variables resolved' message, got: %q", out)
	}
}

func TestFormatResolveSummary_Counts(t *testing.T) {
	results := []ResolveResult{
		{Key: "A", Changed: true},
		{Key: "B", Changed: false, Unresolved: []string{"FOO", "BAR"}},
		{Key: "C", Changed: true},
	}
	summary := FormatResolveSummary(results)
	if !strings.Contains(summary, "2 changed") {
		t.Errorf("expected '2 changed' in summary, got: %q", summary)
	}
	if !strings.Contains(summary, "2 unresolved") {
		t.Errorf("expected '2 unresolved' in summary, got: %q", summary)
	}
}

func TestFormatResolveResults_NoEntries(t *testing.T) {
	out := FormatResolveResults(nil, false)
	if !strings.Contains(out, "no entries") {
		t.Errorf("expected 'no entries' message, got: %q", out)
	}
}
