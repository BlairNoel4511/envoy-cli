package envfile

import (
	"testing"
)

func makeSwapEntries() []Entry {
	return []Entry{
		{Key: "FOO", Value: "foo_val"},
		{Key: "BAR", Value: "bar_val"},
		{Key: "BAZ", Value: "baz_val"},
	}
}

func TestSwap_ExchangesValues(t *testing.T) {
	entries := makeSwapEntries()
	updated, result, summary := Swap(entries, "FOO", "BAR")

	if !result.Swapped {
		t.Fatal("expected Swapped to be true")
	}
	if result.Err != "" {
		t.Fatalf("unexpected error: %s", result.Err)
	}
	if summary.Swapped != 1 || summary.Failed != 0 {
		t.Fatalf("unexpected summary: %+v", summary)
	}

	m := ToMap(updated)
	if m["FOO"] != "bar_val" {
		t.Errorf("expected FOO=bar_val, got %s", m["FOO"])
	}
	if m["BAR"] != "foo_val" {
		t.Errorf("expected BAR=foo_val, got %s", m["BAR"])
	}
	if m["BAZ"] != "baz_val" {
		t.Errorf("expected BAZ unchanged, got %s", m["BAZ"])
	}
}

func TestSwap_KeyANotFound(t *testing.T) {
	entries := makeSwapEntries()
	_, result, summary := Swap(entries, "MISSING", "BAR")

	if result.Swapped {
		t.Fatal("expected Swapped to be false")
	}
	if result.Err == "" {
		t.Fatal("expected an error message")
	}
	if summary.Failed != 1 {
		t.Fatalf("expected 1 failure, got %d", summary.Failed)
	}
}

func TestSwap_KeyBNotFound(t *testing.T) {
	entries := makeSwapEntries()
	_, result, summary := Swap(entries, "FOO", "MISSING")

	if result.Swapped {
		t.Fatal("expected Swapped to be false")
	}
	if result.Err == "" {
		t.Fatal("expected an error message")
	}
	if summary.Failed != 1 {
		t.Fatalf("expected 1 failure, got %d", summary.Failed)
	}
}

func TestSwap_SameKey_ReturnsError(t *testing.T) {
	entries := makeSwapEntries()
	_, result, summary := Swap(entries, "FOO", "FOO")

	if result.Swapped {
		t.Fatal("expected Swapped to be false when swapping key with itself")
	}
	if result.Err == "" {
		t.Fatal("expected an error message for same-key swap")
	}
	if summary.Failed != 1 {
		t.Fatalf("expected 1 failure, got %d", summary.Failed)
	}
}

func TestSwap_RecordsOriginalValues(t *testing.T) {
	entries := makeSwapEntries()
	_, result, _ := Swap(entries, "FOO", "BAZ")

	if result.ValueA != "foo_val" {
		t.Errorf("expected ValueA=foo_val, got %s", result.ValueA)
	}
	if result.ValueB != "baz_val" {
		t.Errorf("expected ValueB=baz_val, got %s", result.ValueB)
	}
}
