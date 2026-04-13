package envfile

import (
	"fmt"
	"strings"
)

// LintRule represents a single lint rule identifier.
type LintRule string

const (
	LintRuleDuplicateKey   LintRule = "duplicate-key"
	LintRuleEmptyValue     LintRule = "empty-value"
	LintRuleNoUppercase    LintRule = "no-uppercase"
	LintRuleTrailingSpace  LintRule = "trailing-space"
	LintRuleSensitiveNoMask LintRule = "sensitive-no-mask"
)

// LintViolation describes a single lint issue found in an entry.
type LintViolation struct {
	Key     string
	Rule    LintRule
	Message string
}

// LintOptions controls which rules are enabled during linting.
type LintOptions struct {
	CheckDuplicateKeys    bool
	CheckEmptyValues      bool
	CheckUppercaseKeys    bool
	CheckTrailingSpace    bool
	CheckSensitiveNoMask  bool
}

// DefaultLintOptions returns a LintOptions with all rules enabled.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		CheckDuplicateKeys:   true,
		CheckEmptyValues:     true,
		CheckUppercaseKeys:   true,
		CheckTrailingSpace:   true,
		CheckSensitiveNoMask: true,
	}
}

// Lint runs all enabled lint rules against the given entries and returns any violations found.
func Lint(entries []Entry, opts LintOptions) []LintViolation {
	var violations []LintViolation
	seen := make(map[string]int)

	for _, e := range entries {
		if opts.CheckDuplicateKeys {
			seen[e.Key]++
			if seen[e.Key] == 2 {
				violations = append(violations, LintViolation{
					Key:     e.Key,
					Rule:    LintRuleDuplicateKey,
					Message: fmt.Sprintf("key %q appears more than once", e.Key),
				})
			}
		}

		if opts.CheckEmptyValues && strings.TrimSpace(e.Value) == "" {
			violations = append(violations, LintViolation{
				Key:     e.Key,
				Rule:    LintRuleEmptyValue,
				Message: fmt.Sprintf("key %q has an empty value", e.Key),
			})
		}

		if opts.CheckUppercaseKeys && e.Key != strings.ToUpper(e.Key) {
			violations = append(violations, LintViolation{
				Key:     e.Key,
				Rule:    LintRuleNoUppercase,
				Message: fmt.Sprintf("key %q is not fully uppercase", e.Key),
			})
		}

		if opts.CheckTrailingSpace && (e.Value != strings.TrimRight(e.Value, " \t")) {
			violations = append(violations, LintViolation{
				Key:     e.Key,
				Rule:    LintRuleTrailingSpace,
				Message: fmt.Sprintf("key %q has trailing whitespace in value", e.Key),
			})
		}

		if opts.CheckSensitiveNoMask && IsSensitive(e.Key) && e.Value == e.Value {
			// Flag sensitive keys that still carry a plaintext-looking value (non-empty, not masked)
			if len(e.Value) > 0 && !strings.Contains(e.Value, "*") {
				violations = append(violations, LintViolation{
					Key:     e.Key,
					Rule:    LintRuleSensitiveNoMask,
					Message: fmt.Sprintf("sensitive key %q appears to have an unmasked value", e.Key),
				})
			}
		}
	}

	return violations
}
