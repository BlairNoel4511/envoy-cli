package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func makeImportEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestImport_AddsNewKeys(t *testing.T) {
	dst := makeImportEntries()
	src := []Entry{{Key: "NEW_KEY", Value: "hello"}}
	out, results, summary := Import(dst, src, ImportOptions{})
	if summary.Added != 1 || summary.Total != 1 {
		t.Fatalf("expected 1 added, got %+v", summary)
	}
	m := ToMap(out)
	if m["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello")
	}
	if results[0].Status != "added" {
		t.Errorf("expected status added, got %s", results[0].Status)
	}
}

func TestImport_SkipsExistingWithoutOverwrite(t *testing.T) {
	dst := makeImportEntries()
	src := []Entry{{Key: "PORT", Value: "9090"}}
	_, _, summary := Import(dst, src, ImportOptions{})
	if summary.Skipped != 1 {
		t.Fatalf("expected 1 skipped, got %+v", summary)
	}
}

func TestImport_OverwritesWhenEnabled(t *testing.T) {
	dst := makeImportEntries()
	src := []Entry{{Key: "PORT", Value: "9090"}}
	out, _, summary := Import(dst, src, ImportOptions{Overwrite: true})
	if summary.Overwritten != 1 {
		t.Fatalf("expected 1 overwritten, got %+v", summary)
	}
	if ToMap(out)["PORT"] != "9090" {
		t.Errorf("expected PORT=9090 after overwrite")
	}
}

func TestImport_SkipsSensitiveKeys(t *testing.T) {
	dst := makeImportEntries()
	src := []Entry{{Key: "SECRET_TOKEN", Value: "abc123"}}
	_, results, summary := Import(dst, src, ImportOptions{SkipSensitive: true})
	if summary.Skipped != 1 {
		t.Fatalf("expected 1 skipped, got %+v", summary)
	}
	if !results[0].Sensitive {
		t.Errorf("expected sensitive flag on result")
	}
}

func TestImport_DryRunDoesNotWrite(t *testing.T) {
	dst := makeImportEntries()
	src := []Entry{{Key: "DRY_KEY", Value: "val"}}
	out, results, _ := Import(dst, src, ImportOptions{DryRun: true})
	if _, ok := ToMap(out)["DRY_KEY"]; ok {
		t.Errorf("dry-run should not write to dst")
	}
	if results[0].Status != "dry-run" {
		t.Errorf("expected dry-run status")
	}
}

func TestImport_FilterByPrefix(t *testing.T) {
	dst := makeImportEntries()
	src := []Entry{
		{Key: "APP_DEBUG", Value: "true"},
		{Key: "DB_HOST", Value: "localhost"},
	}
	_, _, summary := Import(dst, src, ImportOptions{Prefix: "APP_"})
	if summary.Total != 1 {
		t.Fatalf("expected only APP_ prefixed keys, got total=%d", summary.Total)
	}
}

func TestImportFromFile_ParsesEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	_ = os.WriteFile(path, []byte("FOO=bar\nBAZ=qux\n"), 0600)
	entries, err := ImportFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}
