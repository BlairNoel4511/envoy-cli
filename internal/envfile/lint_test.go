package envfile

import (
	"strings"
	"testing"
)

func makeLintEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestLint_NoViolations(t *testing.T) {
	entries := makeLintEntries()
	violations := Lint(entries, DefaultLintOptions())
	// SECRET_KEY will trigger sensitive-no-mask; filter that rule for this test
	var filtered []LintViolation
	for _, v := range violations {
		if v.Rule != LintRuleSensitiveNoMask {
			filtered = append(filtered, v)
		}
	}
	if len(filtered) != 0 {
		t.Errorf("expected no violations, got %d: %+v", len(filtered), filtered)
	}
}

func TestLint_DetectsDuplicateKey(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "FOO", Value: "baz"},
	}
	violations := Lint(entries, LintOptions{CheckDuplicateKeys: true})
	if len(violations) != 1 || violations[0].Rule != LintRuleDuplicateKey {
		t.Errorf("expected duplicate-key violation, got %+v", violations)
	}
}

func TestLint_DetectsEmptyValue(t *testing.T) {
	entries := []Entry{
		{Key: "EMPTY_KEY", Value: ""},
	}
	violations := Lint(entries, LintOptions{CheckEmptyValues: true})
	if len(violations) != 1 || violations[0].Rule != LintRuleEmptyValue {
		t.Errorf("expected empty-value violation, got %+v", violations)
	}
}

func TestLint_DetectsNonUppercaseKey(t *testing.T) {
	entries := []Entry{
		{Key: "mixedCase", Value: "value"},
	}
	violations := Lint(entries, LintOptions{CheckUppercaseKeys: true})
	if len(violations) != 1 || violations[0].Rule != LintRuleNoUppercase {
		t.Errorf("expected no-uppercase violation, got %+v", violations)
	}
}

func TestLint_DetectsTrailingSpace(t *testing.T) {
	entries := []Entry{
		{Key: "SPACED", Value: "value   "},
	}
	violations := Lint(entries, LintOptions{CheckTrailingSpace: true})
	if len(violations) != 1 || violations[0].Rule != LintRuleTrailingSpace {
		t.Errorf("expected trailing-space violation, got %+v", violations)
	}
}

func TestLint_DetectsSensitiveNoMask(t *testing.T) {
	entries := []Entry{
		{Key: "SECRET_TOKEN", Value: "plaintextvalue"},
	}
	violations := Lint(entries, LintOptions{CheckSensitiveNoMask: true})
	if len(violations) != 1 || violations[0].Rule != LintRuleSensitiveNoMask {
		t.Errorf("expected sensitive-no-mask violation, got %+v", violations)
	}
}

func TestFormatLintViolations_NoViolations(t *testing.T) {
	out := FormatLintViolations(nil, false)
	if !strings.Contains(out, "No lint violations") {
		t.Errorf("expected no-violations message, got: %q", out)
	}
}

func TestFormatLintViolations_ContainsRule(t *testing.T) {
	violations := []LintViolation{
		{Key: "FOO", Rule: LintRuleDuplicateKey, Message: `key "FOO" appears more than once`},
	}
	out := FormatLintViolations(violations, false)
	if !strings.Contains(out, "duplicate-key") {
		t.Errorf("expected rule name in output, got: %q", out)
	}
}

func TestFormatLintSummary_PassMessage(t *testing.T) {
	summary := FormatLintSummary(nil)
	if !strings.Contains(summary, "lint passed") {
		t.Errorf("expected pass message, got: %q", summary)
	}
}

func TestFormatLintSummary_FailMessage(t *testing.T) {
	violations := []LintViolation{
		{Key: "FOO", Rule: LintRuleEmptyValue, Message: "empty"},
		{Key: "BAR", Rule: LintRuleEmptyValue, Message: "empty"},
	}
	summary := FormatLintSummary(violations)
	if !strings.Contains(summary, "lint failed") || !strings.Contains(summary, "2") {
		t.Errorf("unexpected summary: %q", summary)
	}
}
