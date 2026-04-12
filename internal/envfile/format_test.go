package envfile

import (
	"strings"
	"testing"
)

func TestFormatDiff_ContainsAddedKey(t *testing.T) {
	entries := []DiffEntry{
		{Key: "NEW_VAR", NewValue: "hello", Type: DiffAdded},
	}
	out := FormatDiff(entries, false)
	if !strings.Contains(out, "+ NEW_VAR=hello") {
		t.Errorf("expected added line, got: %s", out)
	}
}

func TestFormatDiff_ContainsRemovedKey(t *testing.T) {
	entries := []DiffEntry{
		{Key: "OLD_VAR", OldValue: "bye", Type: DiffRemoved},
	}
	out := FormatDiff(entries, false)
	if !strings.Contains(out, "- OLD_VAR=bye") {
		t.Errorf("expected removed line, got: %s", out)
	}
}

func TestFormatDiff_RedactsSensitiveValues(t *testing.T) {
	entries := []DiffEntry{
		{Key: "SECRET_KEY", OldValue: "mysecret", NewValue: "newsecret", Type: DiffChanged},
	}
	out := FormatDiff(entries, false)
	if strings.Contains(out, "mysecret") || strings.Contains(out, "newsecret") {
		t.Errorf("sensitive value should be redacted, got: %s", out)
	}
	if !strings.Contains(out, "****") {
		t.Errorf("expected redaction marker, got: %s", out)
	}
}

func TestFormatDiff_SortedOutput(t *testing.T) {
	entries := []DiffEntry{
		{Key: "ZZZ", NewValue: "z", Type: DiffAdded},
		{Key: "AAA", NewValue: "a", Type: DiffAdded},
	}
	out := FormatDiff(entries, false)
	idxA := strings.Index(out, "AAA")
	idxZ := strings.Index(out, "ZZZ")
	if idxA > idxZ {
		t.Errorf("expected AAA before ZZZ in output")
	}
}

func TestFormatDiff_ColorizeAddsEscapeCodes(t *testing.T) {
	entries := []DiffEntry{
		{Key: "FOO", NewValue: "bar", Type: DiffAdded},
	}
	out := FormatDiff(entries, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes in colorized output, got: %s", out)
	}
}
