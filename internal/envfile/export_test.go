package envfile

import (
	"encoding/json"
	"strings"
	"testing"
)

func makeExportEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "PORT", Value: "8080"},
		{Key: "API_SECRET", Value: "topsecret"},
	}
}

func TestExport_DotEnvFormat(t *testing.T) {
	out, err := Export(makeExportEntries(), FormatDotEnv, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_NAME=envoy") {
		t.Errorf("expected APP_NAME=envoy in output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT=8080 in output, got:\n%s", out)
	}
}

func TestExport_DotEnvRedactsSensitive(t *testing.T) {
	out, err := Export(makeExportEntries(), FormatDotEnv, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("expected DB_PASSWORD to be redacted, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PASSWORD=***") {
		t.Errorf("expected DB_PASSWORD=*** in output, got:\n%s", out)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	out, err := Export(makeExportEntries(), FormatJSON, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, out)
	}
	if m["APP_NAME"] != "envoy" {
		t.Errorf("expected APP_NAME=envoy, got %q", m["APP_NAME"])
	}
}

func TestExport_JSONRedactsSensitive(t *testing.T) {
	out, err := Export(makeExportEntries(), FormatJSON, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["API_SECRET"] != "***" {
		t.Errorf("expected API_SECRET to be redacted, got %q", m["API_SECRET"])
	}
}

func TestExport_ShellFormat(t *testing.T) {
	out, err := Export(makeExportEntries(), FormatShell, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export APP_NAME=") {
		t.Errorf("expected shell export statement, got:\n%s", out)
	}
}

func TestExport_SortedOutput(t *testing.T) {
	out, err := Export(makeExportEntries(), FormatDotEnv, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i := 1; i < len(lines); i++ {
		if lines[i] < lines[i-1] {
			t.Errorf("output is not sorted: %q comes after %q", lines[i], lines[i-1])
		}
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	_, err := Export(makeExportEntries(), ExportFormat("xml"), false)
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}
