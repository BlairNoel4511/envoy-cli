package envfile

import (
	"strings"
	"testing"
)

func TestFormatCommentResults_Empty(t *testing.T) {
	out := FormatCommentResults([]CommentResult{}, false)
	if out != "" {
		t.Errorf("expected empty output for empty results")
	}
}

func TestFormatCommentResults_ShowsAdded(t *testing.T) {
	results := []CommentResult{
		{Key: "APP_NAME", Comment: "the app", Status: "added"},
	}
	out := FormatCommentResults(results, false)
	if !strings.Contains(out, "+ APP_NAME") {
		t.Errorf("expected added marker, got: %s", out)
	}
	if !strings.Contains(out, "the app") {
		t.Errorf("expected comment text in output")
	}
}

func TestFormatCommentResults_ShowsUpdated(t *testing.T) {
	results := []CommentResult{
		{Key: "DB_HOST", Comment: "updated", Status: "updated"},
	}
	out := FormatCommentResults(results, false)
	if !strings.Contains(out, "~ DB_HOST") {
		t.Errorf("expected updated marker, got: %s", out)
	}
}

func TestFormatCommentResults_ShowsNotFound(t *testing.T) {
	results := []CommentResult{
		{Key: "MISSING", Status: "not_found"},
	}
	out := FormatCommentResults(results, false)
	if !strings.Contains(out, "! MISSING") {
		t.Errorf("expected not_found marker, got: %s", out)
	}
}

func TestFormatCommentResults_ColorizeAddsEscapeCodes(t *testing.T) {
	results := []CommentResult{
		{Key: "APP_NAME", Comment: "app", Status: "added"},
	}
	out := FormatCommentResults(results, true)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes in colorized output")
	}
}

func TestFormatCommentSummary_Counts(t *testing.T) {
	sum := CommentSummary{Added: 3, Updated: 1, Unchanged: 2, NotFound: 1}
	out := FormatCommentSummary(sum)
	if !strings.Contains(out, "3 added") {
		t.Errorf("expected added count in summary")
	}
	if !strings.Contains(out, "1 updated") {
		t.Errorf("expected updated count in summary")
	}
	if !strings.Contains(out, "1 not found") {
		t.Errorf("expected not found count in summary")
	}
}
