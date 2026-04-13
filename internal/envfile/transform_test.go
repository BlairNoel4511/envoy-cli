package envfile

import (
	"strings"
	"testing"
)

func makeTransformEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "  production  "},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
		{Key: "DB_PREFIX", Value: "dev_mydb"},
		{Key: "GREETING", Value: "Hello World"},
	}
}

func TestTransform_Uppercase(t *testing.T) {
	entries := makeTransformEntries()
	out, results, err := Transform(entries, TransformOptions{Op: TransformUppercase})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(out))
	}
	for _, r := range results {
		if !r.Changed && r.OldValue != strings.ToUpper(r.OldValue) {
			t.Errorf("key %s should have been changed", r.Key)
		}
	}
	if out[0].Value != "MYAPP" {
		t.Errorf("expected MYAPP, got %s", out[0].Value)
	}
}

func TestTransform_TrimSpace(t *testing.T) {
	entries := makeTransformEntries()
	out, _, err := Transform(entries, TransformOptions{Op: TransformTrimSpace})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "production" {
		t.Errorf("expected trimmed value, got %q", out[1].Value)
	}
}

func TestTransform_SkipSensitive(t *testing.T) {
	entries := makeTransformEntries()
	_, results, err := Transform(entries, TransformOptions{
		Op:            TransformUppercase,
		SkipSensitive: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "SECRET_KEY" && !r.Skipped {
			t.Errorf("expected SECRET_KEY to be skipped")
		}
	}
}

func TestTransform_KeyFilter(t *testing.T) {
	entries := makeTransformEntries()
	_, results, err := Transform(entries, TransformOptions{
		Op:   TransformUppercase,
		Keys: []string{"APP_NAME"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	changedCount := 0
	for _, r := range results {
		if r.Changed {
			changedCount++
		}
	}
	if changedCount != 1 {
		t.Errorf("expected 1 changed, got %d", changedCount)
	}
}

func TestTransform_TrimPrefix(t *testing.T) {
	entries := makeTransformEntries()
	out, _, err := Transform(entries, TransformOptions{
		Op:   TransformTrimPrefix,
		Arg1: "dev_",
		Keys: []string{"DB_PREFIX"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if e.Key == "DB_PREFIX" && e.Value != "mydb" {
			t.Errorf("expected mydb, got %s", e.Value)
		}
	}
}

func TestTransform_Replace(t *testing.T) {
	entries := makeTransformEntries()
	out, _, err := Transform(entries, TransformOptions{
		Op:   TransformReplace,
		Arg1: "World",
		Arg2: "Envoy",
		Keys: []string{"GREETING"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if e.Key == "GREETING" && e.Value != "Hello Envoy" {
			t.Errorf("expected 'Hello Envoy', got %s", e.Value)
		}
	}
}

func TestTransform_UnknownOp(t *testing.T) {
	entries := makeTransformEntries()
	_, _, err := Transform(entries, TransformOptions{Op: "invalid_op"})
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestFormatTransformSummary_Counts(t *testing.T) {
	results := []TransformResult{
		{Key: "A", Changed: true},
		{Key: "B", Changed: true},
		{Key: "C", Skipped: true, Reason: "sensitive key skipped"},
	}
	summary := FormatTransformSummary(results)
	if !strings.Contains(summary, "2 transformed") {
		t.Errorf("expected '2 transformed' in %q", summary)
	}
	if !strings.Contains(summary, "1 skipped") {
		t.Errorf("expected '1 skipped' in %q", summary)
	}
}
