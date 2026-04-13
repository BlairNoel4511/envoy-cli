package envfile

import (
	"strings"
	"testing"
	"time"
)

func TestFormatWatchEvent_Colorize(t *testing.T) {
	ev := WatchEvent{
		Path:      "/app/.env",
		ChangedAt: time.Now(),
		OldSum:    "000000000000",
		NewSum:    "ffffffffffff",
	}
	out := FormatWatchEvent(ev, true)
	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI escape codes when colorize=true")
	}
}

func TestFormatWatchEvent_NoColor(t *testing.T) {
	ev := WatchEvent{
		Path:      "/app/.env",
		ChangedAt: time.Now(),
		OldSum:    "000000000000",
		NewSum:    "ffffffffffff",
	}
	out := FormatWatchEvent(ev, false)
	if strings.Contains(out, "\033[") {
		t.Error("expected no ANSI escape codes when colorize=false")
	}
}

func TestShortSum_TruncatesLongSum(t *testing.T) {
	long := "abcdef1234567890abcdef"
	result := shortSum(long)
	if len(result) >= len(long) {
		t.Errorf("expected truncated sum, got %q", result)
	}
	if !strings.HasSuffix(result, "...") {
		t.Errorf("expected ellipsis suffix, got %q", result)
	}
}

func TestShortSum_ShortSumUnchanged(t *testing.T) {
	short := "abc123"
	result := shortSum(short)
	if result != short {
		t.Errorf("expected %q unchanged, got %q", short, result)
	}
}

func TestFormatWatchSummary_ZeroCount(t *testing.T) {
	s := FormatWatchSummary("prod.env", 0)
	if !strings.Contains(s, "no changes") {
		t.Errorf("expected 'no changes' in summary, got: %s", s)
	}
	if !strings.Contains(s, "prod.env") {
		t.Errorf("expected filename in summary, got: %s", s)
	}
}

func TestFormatWatchSummary_MultipleChanges(t *testing.T) {
	s := FormatWatchSummary("staging.env", 5)
	if !strings.Contains(s, "5 changes") {
		t.Errorf("expected '5 changes', got: %s", s)
	}
}
