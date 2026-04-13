package envfile

import (
	"testing"
)

func makePatchEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "development"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestPatch_SetNewKey(t *testing.T) {
	entries := makePatchEntries()
	ins := []PatchInstruction{{Op: PatchSet, Key: "NEW_KEY", Value: "hello"}}
	out, results := Patch(entries, ins, PatchOptions{})
	if !results[0].Applied {
		t.Fatal("expected applied")
	}
	_, ok := Lookup(out, "NEW_KEY")
	if !ok {
		t.Fatal("expected NEW_KEY in output")
	}
}

func TestPatch_SetExistingWithoutOverwrite(t *testing.T) {
	entries := makePatchEntries()
	ins := []PatchInstruction{{Op: PatchSet, Key: "APP_ENV", Value: "production"}}
	_, results := Patch(entries, ins, PatchOptions{Overwrite: false})
	if !results[0].Skipped {
		t.Fatal("expected skipped")
	}
}

func TestPatch_SetExistingWithOverwrite(t *testing.T) {
	entries := makePatchEntries()
	ins := []PatchInstruction{{Op: PatchSet, Key: "APP_ENV", Value: "production"}}
	out, results := Patch(entries, ins, PatchOptions{Overwrite: true})
	if !results[0].Applied {
		t.Fatal("expected applied")
	}
	v, _ := Lookup(out, "APP_ENV")
	if v != "production" {
		t.Fatalf("expected production, got %s", v)
	}
}

func TestPatch_SetIdenticalValueSkips(t *testing.T) {
	entries := makePatchEntries()
	ins := []PatchInstruction{{Op: PatchSet, Key: "APP_ENV", Value: "development"}}
	_, results := Patch(entries, ins, PatchOptions{Overwrite: true})
	if !results[0].Skipped {
		t.Fatal("expected skipped for identical value")
	}
}

func TestPatch_DeleteExistingKey(t *testing.T) {
	entries := makePatchEntries()
	ins := []PatchInstruction{{Op: PatchDelete, Key: "DB_HOST"}}
	out, results := Patch(entries, ins, PatchOptions{})
	if !results[0].Applied {
		t.Fatal("expected applied")
	}
	_, ok := Lookup(out, "DB_HOST")
	if ok {
		t.Fatal("expected DB_HOST to be removed")
	}
}

func TestPatch_DeleteMissingKey(t *testing.T) {
	entries := makePatchEntries()
	ins := []PatchInstruction{{Op: PatchDelete, Key: "MISSING"}}
	_, results := Patch(entries, ins, PatchOptions{})
	if !results[0].Skipped {
		t.Fatal("expected skipped")
	}
}

func TestPatch_RenameKey(t *testing.T) {
	entries := makePatchEntries()
	ins := []PatchInstruction{{Op: PatchRename, Key: "APP_ENV", NewKey: "APP_ENVIRONMENT"}}
	out, results := Patch(entries, ins, PatchOptions{})
	if !results[0].Applied {
		t.Fatal("expected applied")
	}
	_, ok := Lookup(out, "APP_ENVIRONMENT")
	if !ok {
		t.Fatal("expected APP_ENVIRONMENT in output")
	}
}

func TestPatch_DryRunDoesNotMutate(t *testing.T) {
	entries := makePatchEntries()
	ins := []PatchInstruction{{Op: PatchSet, Key: "NEW_KEY", Value: "val"}}
	out, results := Patch(entries, ins, PatchOptions{DryRun: true})
	if !results[0].Applied {
		t.Fatal("expected applied in dry run")
	}
	_, ok := Lookup(out, "NEW_KEY")
	if ok {
		t.Fatal("dry run should not mutate entries")
	}
}

func TestPatch_UnknownOpSkipped(t *testing.T) {
	entries := makePatchEntries()
	ins := []PatchInstruction{{Op: PatchOp("noop"), Key: "APP_ENV"}}
	_, results := Patch(entries, ins, PatchOptions{})
	if !results[0].Skipped {
		t.Fatal("expected skipped for unknown op")
	}
}
