package envfile

import (
	"os"
	"testing"
	"time"
)

func writeTempWatchFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "watch_*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestWatch_DetectsChange(t *testing.T) {
	path := writeTempWatchFile(t, "KEY=original\n")

	done := make(chan struct{})
	defer close(done)

	opts := WatchOptions{Interval: 50 * time.Millisecond, MaxChecks: 20}
	events, err := Watch(path, opts, done)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	// Modify the file after a short delay.
	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = os.WriteFile(path, []byte("KEY=modified\n"), 0o600)
	}()

	select {
	case ev, ok := <-events:
		if !ok {
			t.Fatal("events channel closed before receiving event")
		}
		if ev.Path != path {
			t.Errorf("expected path %q, got %q", path, ev.Path)
		}
		if ev.OldSum == ev.NewSum {
			t.Error("expected old and new checksums to differ")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for watch event")
	}
}

func TestWatch_NoEventWhenUnchanged(t *testing.T) {
	path := writeTempWatchFile(t, "KEY=stable\n")

	done := make(chan struct{})

	opts := WatchOptions{Interval: 30 * time.Millisecond, MaxChecks: 5}
	events, err := Watch(path, opts, done)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	// Wait for the watcher to exhaust its checks.
	time.Sleep(250 * time.Millisecond)
	close(done)

	count := 0
	for range events {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 events for unchanged file, got %d", count)
	}
}

func TestWatch_FileNotFound_ReturnsError(t *testing.T) {
	_, err := Watch("/nonexistent/path/.env", DefaultWatchOptions(), make(chan struct{}))
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestFormatWatchEvent_ContainsPath(t *testing.T) {
	ev := WatchEvent{
		Path:      "/app/.env",
		ChangedAt: time.Now(),
		OldSum:    "aabbcc112233",
		NewSum:    "ddeeff445566",
	}
	out := FormatWatchEvent(ev, false)
	if !containsStr(out, "/app/.env") {
		t.Errorf("expected path in output, got:\n%s", out)
	}
	if !containsStr(out, "aabbcc112233") {
		t.Errorf("expected old sum in output, got:\n%s", out)
	}
}

func TestFormatWatchSummary_PluralAndSingular(t *testing.T) {
	if s := FormatWatchSummary(".env", 0); !containsStr(s, "no changes") {
		t.Errorf("unexpected zero summary: %s", s)
	}
	if s := FormatWatchSummary(".env", 1); !containsStr(s, "1 change") {
		t.Errorf("unexpected singular summary: %s", s)
	}
	if s := FormatWatchSummary(".env", 3); !containsStr(s, "3 changes") {
		t.Errorf("unexpected plural summary: %s", s)
	}
}

// containsStr is a small helper used across watch tests.
func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
