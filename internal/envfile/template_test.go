package envfile

import (
	"strings"
	"testing"
)

func makeTemplateEntries() []Entry {
	return []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "5432"},
		{Key: "DB_URL", Value: "postgres://${HOST}:${PORT}/mydb"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
	}
}

func TestExpandTemplate_ResolvesPlaceholders(t *testing.T) {
	entries := makeTemplateEntries()
	r := ExpandTemplate("connect to ${HOST}:${PORT}", entries)
	if r.Expanded != "connect to localhost:5432" {
		t.Errorf("unexpected expansion: %q", r.Expanded)
	}
	if len(r.Missing) != 0 {
		t.Errorf("expected no missing, got %v", r.Missing)
	}
	if len(r.Resolved) != 2 {
		t.Errorf("expected 2 resolved, got %d", len(r.Resolved))
	}
}

func TestExpandTemplate_TracksMissingVars(t *testing.T) {
	entries := makeTemplateEntries()
	r := ExpandTemplate("${HOST}:${UNDEFINED}", entries)
	if !strings.Contains(r.Expanded, "localhost") {
		t.Errorf("expected HOST to be resolved")
	}
	if len(r.Missing) != 1 || r.Missing[0] != "UNDEFINED" {
		t.Errorf("expected UNDEFINED in missing, got %v", r.Missing)
	}
}

func TestExpandTemplate_NoDollarSign_Unchanged(t *testing.T) {
	entries := makeTemplateEntries()
	r := ExpandTemplate("no variables here", entries)
	if r.Expanded != "no variables here" {
		t.Errorf("expected unchanged string, got %q", r.Expanded)
	}
	if len(r.Resolved) != 0 || len(r.Missing) != 0 {
		t.Error("expected empty resolved and missing")
	}
}

func TestExpandTemplate_DeduplicatesMissing(t *testing.T) {
	entries := makeTemplateEntries()
	r := ExpandTemplate("${MISSING} and ${MISSING} again", entries)
	if len(r.Missing) != 1 {
		t.Errorf("expected 1 unique missing entry, got %d", len(r.Missing))
	}
}

func TestExpandTemplate_EmptyTemplate(t *testing.T) {
	entries := makeTemplateEntries()
	r := ExpandTemplate("", entries)
	if r.Expanded != "" {
		t.Errorf("expected empty string, got %q", r.Expanded)
	}
	if len(r.Resolved) != 0 || len(r.Missing) != 0 {
		t.Error("expected empty resolved and missing for empty template")
	}
}

func TestExpandEntries_ExpandsValues(t *testing.T) {
	entries := makeTemplateEntries()
	expanded, missing := ExpandEntries(entries)

	m := ToMap(expanded)
	if m["DB_URL"] != "postgres://localhost:5432/mydb" {
		t.Errorf("unexpected DB_URL: %q", m["DB_URL"])
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing vars, got %v", missing)
	}
}

func TestExpandEntries_ReportsMissingVars(t *testing.T) {
	entries := []Entry{
		{Key: "URL", Value: "http://${UNDEFINED_HOST}/path"},
	}
	_, missing := ExpandEntries(entries)
	if len(missing) == 0 {
		t.Error("expected missing vars to be reported")
	}
}

func TestFormatTemplateResult_ContainsExpanded(t *testing.T) {
	r := TemplateResult{
		Expanded: "hello world",
		Resolved: []string{"GREETING"},
		Missing:  []string{"NAME"},
	}
	out := FormatTemplateResult(r)
	if !strings.Contains(out, "hello world") {
		t.Error("expected expanded value in output")
	}
	if !strings.Contains(out, "GREETING") {
		t.Error("expected resolved key in output")
	}
	if !strings.Contains(out, "NAME") {
		t.Error("expected missing key in output")
	}
}
