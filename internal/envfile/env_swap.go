package envfile

import "fmt"

// SwapResult describes the outcome of swapping two keys' values.
type SwapResult struct {
	KeyA    string
	KeyB    string
	ValueA  string // original value of KeyA (now assigned to KeyB)
	ValueB  string // original value of KeyB (now assigned to KeyA)
	Swapped bool
	Err     string
}

// SwapSummary holds aggregate counts for a swap operation.
type SwapSummary struct {
	Swapped int
	Failed  int
}

// Swap exchanges the values of two keys within the given entries.
// It returns the updated entries, a SwapResult, and a SwapSummary.
func Swap(entries []Entry, keyA, keyB string) ([]Entry, SwapResult, SwapSummary) {
	result := SwapResult{KeyA: keyA, KeyB: keyB}

	idxA := -1
	idxB := -1
	for i, e := range entries {
		if e.Key == keyA {
			idxA = i
		}
		if e.Key == keyB {
			idxB = i
		}
	}

	if idxA == -1 {
		result.Err = fmt.Sprintf("key %q not found", keyA)
		return entries, result, SwapSummary{Failed: 1}
	}
	if idxB == -1 {
		result.Err = fmt.Sprintf("key %q not found", keyB)
		return entries, result, SwapSummary{Failed: 1}
	}
	if keyA == keyB {
		result.Err = "cannot swap a key with itself"
		return entries, result, SwapSummary{Failed: 1}
	}

	result.ValueA = entries[idxA].Value
	result.ValueB = entries[idxB].Value

	entries[idxA].Value = result.ValueB
	entries[idxB].Value = result.ValueA
	result.Swapped = true

	return entries, result, SwapSummary{Swapped: 1}
}
