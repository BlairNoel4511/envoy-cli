package envfile

import (
	"testing"
)

func makeCommentEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "localhost", Comment: "existing comment"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
	}
}

func TestComment_AddsCommentToKey(t *testing.T) {
	entries := makeCommentEntries()
	result, res := Comment(entries, "APP_NAME", "the application name", CommentOptions{})
	if res.Status != "added" {
		t.Fatalf("expected added, got %s", res.Status)
	}
	if result[0].Comment != "the application name" {
		t.Errorf("expected comment to be set")
	}
}

func TestComment_UpdatesExistingWithOverwrite(t *testing.T) {
	entries := makeCommentEntries()
	_, res := Comment(entries, "DB_HOST", "new comment", CommentOptions{Overwrite: true})
	if res.Status != "updated" {
		t.Fatalf("expected updated, got %s", res.Status)
	}
	if res.Comment != "new comment" {
		t.Errorf("expected new comment, got %s", res.Comment)
	}
}

func TestComment_SkipsExistingWithoutOverwrite(t *testing.T) {
	entries := makeCommentEntries()
	_, res := Comment(entries, "DB_HOST", "new comment", CommentOptions{Overwrite: false})
	if res.Status != "unchanged" {
		t.Fatalf("expected unchanged, got %s", res.Status)
	}
	if res.Comment != "existing comment" {
		t.Errorf("expected original comment preserved")
	}
}

func TestComment_KeyNotFound(t *testing.T) {
	entries := makeCommentEntries()
	_, res := Comment(entries, "MISSING_KEY", "comment", CommentOptions{})
	if res.Status != "not_found" {
		t.Fatalf("expected not_found, got %s", res.Status)
	}
}

func TestComment_UnchangedWhenIdentical(t *testing.T) {
	entries := makeCommentEntries()
	_, res := Comment(entries, "DB_HOST", "existing comment", CommentOptions{Overwrite: true})
	if res.Status != "unchanged" {
		t.Fatalf("expected unchanged, got %s", res.Status)
	}
}

func TestCommentMany_SummaryCountsCorrect(t *testing.T) {
	entries := makeCommentEntries()
	comments := map[string]string{
		"APP_NAME":   "app name",
		"DB_HOST":    "db host",
		"MISSING":    "nope",
		"SECRET_KEY": "secret",
	}
	_, _, sum := CommentMany(entries, comments, CommentOptions{Overwrite: false})
	if sum.Added != 2 {
		t.Errorf("expected 2 added, got %d", sum.Added)
	}
	if sum.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", sum.Unchanged)
	}
	if sum.NotFound != 1 {
		t.Errorf("expected 1 not_found, got %d", sum.NotFound)
	}
}
