package envfile

import (
	"strings"
	"testing"
)

func makeMaskEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "supersecret123"},
		{Key: "API_KEY", Value: "abcdef1234567890"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestMaskValue_BasicMasking(t *testing.T) {
	opts := DefaultMaskOptions()
	result := MaskValue("supersecret123", opts)
	if result == "supersecret123" {
		t.Error("expected value to be masked")
	}
	if !strings.Contains(result, "**") {
		t.Errorf("expected mask chars in result, got %q", result)
	}
}

func TestMaskValue_PreservesPrefix(t *testing.T) {
	opts := DefaultMaskOptions()
	opts.VisibleChars = 3
	result := MaskValue("abcXXXXXXXyz", opts)
	if !strings.HasPrefix(result, "abc") {
		t.Errorf("expected prefix 'abc', got %q", result)
	}
}

func TestMaskValue_EmptyValue(t *testing.T) {
	opts := DefaultMaskOptions()
	result := MaskValue("", opts)
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestMaskValue_ShortValue(t *testing.T) {
	opts := DefaultMaskOptions()
	opts.VisibleChars = 4
	opts.MaskLength = 4
	result := MaskValue("ab", opts)
	if result == "ab" {
		t.Error("expected short value to still be masked")
	}
}

func TestMaskEntries_SensitiveOnlySkipsPlain(t *testing.T) {
	entries := makeMaskEntries()
	opts := DefaultMaskOptions()
	opts.SensitiveOnly = true
	results := MaskEntries(entries, opts)
	for _, r := range results {
		if r.Key == "APP_NAME" && r.WasMasked {
			t.Errorf("APP_NAME should not be masked")
		}
		if r.Key == "PORT" && r.WasMasked {
			t.Errorf("PORT should not be masked")
		}
	}
}

func TestMaskEntries_MasksSensitiveKeys(t *testing.T) {
	entries := makeMaskEntries()
	opts := DefaultMaskOptions()
	results := MaskEntries(entries, opts)
	for _, r := range results {
		if r.Key == "DB_PASSWORD" && !r.WasMasked {
			t.Errorf("DB_PASSWORD should be masked")
		}
		if r.Key == "API_KEY" && !r.WasMasked {
			t.Errorf("API_KEY should be masked")
		}
	}
}

func TestFormatMaskResults_ContainsKeys(t *testing.T) {
	entries := makeMaskEntries()
	opts := DefaultMaskOptions()
	results := MaskEntries(entries, opts)
	out := FormatMaskResults(results, false)
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD in output, got:\n%s", out)
	}
}

func TestFormatMaskResults_ColorizeAddsEscapeCodes(t *testing.T) {
	entries := []Entry{{Key: "SECRET_KEY", Value: "topsecret"}}
	opts := DefaultMaskOptions()
	results := MaskEntries(entries, opts)
	out := FormatMaskResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes in colorized output, got:\n%s", out)
	}
}

func TestFormatMaskSummary_Counts(t *testing.T) {
	entries := makeMaskEntries()
	opts := DefaultMaskOptions()
	results := MaskEntries(entries, opts)
	summary := FormatMaskSummary(results)
	if !strings.Contains(summary, "masked") {
		t.Errorf("expected 'masked' in summary, got %q", summary)
	}
}
