package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSync_WritesNewKeys(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, ".env")

	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := Sync(dest, src, SyncOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Written) != 2 {
		t.Errorf("expected 2 written, got %d", len(res.Written))
	}
}

func TestSync_SkipsExistingWithoutOverwrite(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, ".env")
	_ = os.WriteFile(dest, []byte("FOO=original\n"), 0o644)

	src := map[string]string{"FOO": "changed"}
	res, err := Sync(dest, src, SyncOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestSync_OverwritesExistingWhenEnabled(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, ".env")
	_ = os.WriteFile(dest, []byte("FOO=original\n"), 0o644)

	src := map[string]string{"FOO": "changed"}
	res, err := Sync(dest, src, SyncOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten, got %d", len(res.Overwritten))
	}
}

func TestSync_DryRunDoesNotWrite(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, ".env")

	src := map[string]string{"FOO": "bar"}
	_, err := Sync(dest, src, SyncOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
		t.Error("expected destination file to not exist after dry run")
	}
}

func TestSync_SkipsIdenticalValues(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, ".env")
	_ = os.WriteFile(dest, []byte("FOO=bar\n"), 0o644)

	src := map[string]string{"FOO": "bar"}
	res, err := Sync(dest, src, SyncOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped (identical value), got %d", len(res.Skipped))
	}
}
