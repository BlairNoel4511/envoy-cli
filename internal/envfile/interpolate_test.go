package envfile

import "testing"

func TestInterpolate_BraceStyle(t *testing.T) {
	result, missing := interpolate("${FOO} bar", map[string]string{"FOO": "hello"})
	if result != "hello bar" {
		t.Errorf("expected 'hello bar', got %q", result)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing vars, got %v", missing)
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	result, missing := interpolate("$FOO world", map[string]string{"FOO": "hi"})
	if result != "hi world" {
		t.Errorf("expected 'hi world', got %q", result)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing vars, got %v", missing)
	}
}

func TestInterpolate_MissingVar(t *testing.T) {
	_, missing := interpolate("${NOPE}", map[string]string{})
	if len(missing) != 1 || missing[0] != "NOPE" {
		t.Errorf("expected [NOPE] in missing, got %v", missing)
	}
}

func TestInterpolate_DedupMissing(t *testing.T) {
	_, missing := interpolate("${X} and ${X}", map[string]string{})
	if len(missing) != 1 {
		t.Errorf("expected deduped missing, got %v", missing)
	}
}

func TestInterpolate_NoVars(t *testing.T) {
	result, missing := interpolate("plain value", map[string]string{})
	if result != "plain value" {
		t.Errorf("expected unchanged, got %q", result)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing, got %v", missing)
	}
}

func TestInterpolate_MalformedBrace(t *testing.T) {
	result, _ := interpolate("${", map[string]string{})
	if result != "${" {
		t.Errorf("expected malformed brace preserved, got %q", result)
	}
}

func TestInterpolate_LonelyDollar(t *testing.T) {
	result, _ := interpolate("price is $", map[string]string{})
	if result != "price is $" {
		t.Errorf("expected lonely dollar preserved, got %q", result)
	}
}
