package envfile

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	tmp, err := os.CreateTemp(t.TempDir(), ".env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := tmp.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	tmp.Close()
	return tmp.Name()
}

func TestParse_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDEBUG=false\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := env.ToMap()
	if m["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", m["APP_ENV"])
	}
	if m["DEBUG"] != "false" {
		t.Errorf("expected DEBUG=false, got %q", m["DEBUG"])
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"` + "\n" + `TOKEN='abc123'` + "\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := env.ToMap()
	if m["SECRET"] != "my secret value" {
		t.Errorf("expected unquoted value, got %q", m["SECRET"])
	}
	if m["TOKEN"] != "abc123" {
		t.Errorf("expected unquoted value, got %q", m["TOKEN"])
	}
}

func TestParse_CommentsIgnoredInMap(t *testing.T) {
	path := writeTempEnv(t, "# This is a comment\nKEY=value\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := env.ToMap()
	if _, ok := m[""]; ok {
		t.Error("comment should not appear as empty key in map")
	}
	if m["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", m["KEY"])
	}
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE_NO_EQUALS\n")
	_, err := Parse(path)
	if err == nil {
		t.Error("expected error for invalid line, got nil")
	}
}

func TestParse_FileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
