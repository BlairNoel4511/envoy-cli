package envfile

import (
	"strings"
	"testing"
)

func TestFormatGroupList_Empty(t *testing.T) {
	gs := NewGroupStore()
	out := FormatGroupList(gs)
	if !strings.Contains(out, "No groups") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatGroupList_ContainsGroupNames(t *testing.T) {
	gs := NewGroupStore()
	gs.Add("database", []string{"DB_HOST", "DB_PASS"})
	gs.Add("app", []string{"APP_PORT"})
	out := FormatGroupList(gs)
	if !strings.Contains(out, "database") {
		t.Errorf("expected 'database' in output: %s", out)
	}
	if !strings.Contains(out, "app") {
		t.Errorf("expected 'app' in output: %s", out)
	}
}

func TestFormatGroupDetail_NilGroup(t *testing.T) {
	out := FormatGroupDetail(nil)
	if !strings.Contains(out, "not found") {
		t.Errorf("expected not-found message, got: %s", out)
	}
}

func TestFormatGroupDetail_ShowsKeys(t *testing.T) {
	g := &Group{Name: "infra", Keys: []string{"DB_HOST", "APP_PORT"}}
	out := FormatGroupDetail(g)
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output: %s", out)
	}
	if !strings.Contains(out, "APP_PORT") {
		t.Errorf("expected APP_PORT in output: %s", out)
	}
}

func TestFormatFilteredByGroup_RedactsSensitive(t *testing.T) {
	entries := []Entry{
		{Key: "DB_PASSWORD", Value: "supersecret"},
		{Key: "APP_PORT", Value: "8080"},
	}
	out := FormatFilteredByGroup(entries, "mixed")
	if strings.Contains(out, "supersecret") {
		t.Error("sensitive value should be redacted")
	}
	if !strings.Contains(out, "***") {
		t.Error("expected redaction marker ***")
	}
	if !strings.Contains(out, "8080") {
		t.Error("expected non-sensitive value to appear")
	}
}

func TestFormatFilteredByGroup_Empty(t *testing.T) {
	out := FormatFilteredByGroup([]Entry{}, "empty-group")
	if !strings.Contains(out, "No entries") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
