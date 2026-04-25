package envfile

import (
	"testing"
)

func makeFlattenEntries() []Entry {
	return []Entry{
		{Key: "APP__HOST", Value: "localhost"},
		{Key: "APP__PORT", Value: "8080"},
		{Key: "DB__PASSWORD", Value: "secret"},
		{Key: "CLEAN_KEY", Value: "unchanged"},
		{Key: "__LEADING", Value: "trim"},
	}
}

func TestFlatten_CollapsesDoubleSeparators(t *testing.T) {
	entries := makeFlattenEntries()
	out, results, summary := Flatten(entries, FlattenOptions{})

	if summary.Total != 5 {
		t.Fatalf("expected 5 total, got %d", summary.Total)
	}
	if summary.Changed == 0 {
		t.Error("expected at least one changed key")
	}
	if out[0].Key != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %s", out[0].Key)
	}
	if results[0].Changed != true {
		t.Error("expected APP__HOST to be marked changed")
	}
}

func TestFlatten_UppercaseOption(t *testing.T) {
	entries := []Entry{
		{Key: "app__name", Value: "envoy"},
	}
	out, _, _ := Flatten(entries, FlattenOptions{Uppercase: true})
	if out[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME, got %s", out[0].Key)
	}
}

func TestFlatten_UnchangedKeyCountedAsSkipped(t *testing.T) {
	entries := []Entry{
		{Key: "ALREADY_CLEAN", Value: "ok"},
	}
	_, _, summary := Flatten(entries, FlattenOptions{})
	if summary.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", summary.Skipped)
	}
	if summary.Changed != 0 {
		t.Errorf("expected 0 changed, got %d", summary.Changed)
	}
}

func TestFlatten_RedactsSensitiveValues(t *testing.T) {
	entries := []Entry{
		{Key: "DB__PASSWORD", Value: "supersecret"},
	}
	_, results, _ := Flatten(entries, FlattenOptions{Redact: true})
	if results[0].Value != "***" {
		t.Errorf("expected redacted value, got %s", results[0].Value)
	}
	if !results[0].Sensitive {
		t.Error("expected Sensitive to be true for PASSWORD key")
	}
}

func TestFlatten_StripLeadingAndTrailingSeparator(t *testing.T) {
	entries := []Entry{
		{Key: "__LEADING", Value: "trim"},
	}
	out, _, _ := Flatten(entries, FlattenOptions{})
	if out[0].Key != "LEADING" {
		t.Errorf("expected LEADING, got %s", out[0].Key)
	}
}

func TestFlatten_PreservesOriginalValueInOutput(t *testing.T) {
	entries := []Entry{
		{Key: "DB__PASSWORD", Value: "secret123"},
	}
	out, results, _ := Flatten(entries, FlattenOptions{Redact: true})
	// Output entries should contain real value.
	if out[0].Value != "secret123" {
		t.Errorf("expected real value in output entry, got %s", out[0].Value)
	}
	// Results should contain masked value.
	if results[0].Value != "***" {
		t.Errorf("expected masked value in result, got %s", results[0].Value)
	}
}

func TestFlatten_SummaryTotalMatchesInputLength(t *testing.T) {
	entries := makeFlattenEntries()
	_, _, summary := Flatten(entries, FlattenOptions{})
	if summary.Total != len(entries) {
		t.Errorf("expected summary.Total %d to match input length %d", summary.Total, len(entries))
	}
	if summary.Changed+summary.Skipped != summary.Total {
		t.Errorf("expected Changed (%d) + Skipped (%d) to equal Total (%d)",
			summary.Changed, summary.Skipped, summary.Total)
	}
}
