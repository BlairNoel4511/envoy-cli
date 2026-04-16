package envfile

import (
	"testing"
)

func makeCopyEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "SECRET_KEY", Value: "supersecret"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestCopy_CopiesKeyToNewDest(t *testing.T) {
	entries, result := Copy(makeCopyEntries(), "APP_NAME", "APP_NAME_COPY", CopyOptions{})
	if result.Status != "copied" {
		t.Fatalf("expected copied, got %s", result.Status)
	}
	v, ok := Lookup(entries, "APP_NAME_COPY")
	if !ok || v.Value != "myapp" {
		t.Errorf("expected APP_NAME_COPY=myapp")
	}
}

func TestCopy_SkipsExistingWithoutOverwrite(t *testing.T) {
	_, result := Copy(makeCopyEntries(), "APP_NAME", "PORT", CopyOptions{Overwrite: false})
	if result.Status != "skipped" {
		t.Errorf("expected skipped, got %s", result.Status)
	}
}

func TestCopy_OverwritesExistingWhenEnabled(t *testing.T) {
	entries, result := Copy(makeCopyEntries(), "APP_NAME", "PORT", CopyOptions{Overwrite: true})
	if result.Status != "overwritten" {
		t.Fatalf("expected overwritten, got %s", result.Status)
	}
	v, _ := Lookup(entries, "PORT")
	if v.Value != "myapp" {
		t.Errorf("expected PORT=myapp after overwrite")
	}
}

func TestCopy_SourceNotFound(t *testing.T) {
	_, result := Copy(makeCopyEntries(), "MISSING", "DEST", CopyOptions{})
	if result.Status != "not_found" {
		t.Errorf("expected not_found, got %s", result.Status)
	}
}

func TestCopy_RedactsSensitiveValue(t *testing.T) {
	_, result := Copy(makeCopyEntries(), "SECRET_KEY", "SECRET_COPY", CopyOptions{RedactSensitive: true})
	if result.Value != "***" {
		t.Errorf("expected redacted value, got %s", result.Value)
	}
}

func TestCopySummary_Messages(t *testing.T) {
	cases := []struct {
		status string
		contains string
	}{
		{"copied", "copied"},
		{"overwritten", "overwritten"},
		{"skipped", "skipped"},
		{"not_found", "error"},
	}
	for _, c := range cases {
		r := CopyResult{SourceKey: "A", DestKey: "B", Status: c.status}
		s := CopySummary(r)
		if len(s) == 0 {
			t.Errorf("expected non-empty summary for status %s", c.status)
		}
	}
}
