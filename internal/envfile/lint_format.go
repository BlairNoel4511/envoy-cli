package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatLintViolations returns a human-readable report of lint violations.
// If colorize is true, ANSI escape codes are added for terminal output.
func FormatLintViolations(violations []LintViolation, colorize bool) string {
	if len(violations) == 0 {
		if colorize {
			return "\033[32m✔ No lint violations found.\033[0m\n"
		}
		return "✔ No lint violations found.\n"
	}

	sorted := make([]LintViolation, len(violations))
	copy(sorted, violations)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Key != sorted[j].Key {
			return sorted[i].Key < sorted[j].Key
		}
		return sorted[i].Rule < sorted[j].Rule
	})

	var sb strings.Builder
	for _, v := range sorted {
		line := fmt.Sprintf("  [%s] %s\n", v.Rule, v.Message)
		if colorize {
			line = fmt.Sprintf("  \033[33m[%s]\033[0m %s\n", v.Rule, v.Message)
		}
		sb.WriteString(line)
	}
	return sb.String()
}

// FormatLintSummary returns a one-line summary of lint results.
func FormatLintSummary(violations []LintViolation) string {
	if len(violations) == 0 {
		return "lint passed: 0 violations"
	}

	counts := make(map[LintRule]int)
	for _, v := range violations {
		counts[v.Rule]++
	}

	parts := make([]string, 0, len(counts))
	for rule, n := range counts {
		parts = append(parts, fmt.Sprintf("%s=%d", rule, n))
	}
	sort.Strings(parts)

	return fmt.Sprintf("lint failed: %d violation(s) [%s]", len(violations), strings.Join(parts, ", "))
}

// ViolationRules returns the unique set of rules triggered across all violations.
func ViolationRules(violations []LintViolation) []LintRule {
	seen := make(map[LintRule]struct{})
	var rules []LintRule
	for _, v := range violations {
		if _, ok := seen[v.Rule]; !ok {
			seen[v.Rule] = struct{}{}
			rules = append(rules, v.Rule)
		}
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i] < rules[j] })
	return rules
}
