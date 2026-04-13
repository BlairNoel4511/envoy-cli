package envfile

import (
	"strings"
	"testing"
)

func TestFormatTagList_Empty(t *testing.T) {
	ts := NewTagStore()
	out := FormatTagList(ts)
	if out != "(no tags)" {
		t.Errorf("expected '(no tags)', got %q", out)
	}
}

func TestFormatTagList_ContainsKeys(t *testing.T) {
	ts := NewTagStore()
	ts.Add("DB_HOST", "infra")
	ts.Add("API_KEY", "secret")

	out := FormatTagList(ts)
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Error("expected API_KEY in output")
	}
}

func TestFormatTagsForKey_NoTags(t *testing.T) {
	out := FormatTagsForKey("PORT", nil)
	if !strings.Contains(out, "no tags") {
		t.Errorf("expected 'no tags' message, got %q", out)
	}
}

func TestFormatTagsForKey_WithLabels(t *testing.T) {
	out := FormatTagsForKey("DB_HOST", []string{"infra", "required"})
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected key name in output")
	}
	if !strings.Contains(out, "infra") || !strings.Contains(out, "required") {
		t.Error("expected labels in output")
	}
}

func TestFormatKeysWithTag_NoKeys(t *testing.T) {
	out := FormatKeysWithTag("infra", nil)
	if !strings.Contains(out, "No keys") {
		t.Errorf("expected 'No keys' message, got %q", out)
	}
}

func TestFormatKeysWithTag_WithKeys(t *testing.T) {
	out := FormatKeysWithTag("infra", []string{"DB_HOST", "DB_PORT"})
	if !strings.Contains(out, "infra") {
		t.Error("expected label name in output")
	}
	if !strings.Contains(out, "DB_HOST") || !strings.Contains(out, "DB_PORT") {
		t.Error("expected key names in output")
	}
}
