package envfile

import (
	"fmt"
	"time"
)

// RollbackEntry represents a single key-value change that can be undone.
type RollbackEntry struct {
	Key      string
	OldValue string
	NewValue string
	HadKey   bool // false if the key didn't exist before the change
}

// RollbackPlan holds a set of reversible changes for a given operation.
type RollbackPlan struct {
	OperationID string
	CreatedAt   time.Time
	Entries     []RollbackEntry
}

// NewRollbackPlan creates a RollbackPlan by comparing before and after entry slices.
func NewRollbackPlan(operationID string, before, after []Entry) RollbackPlan {
	beforeMap := ToMap(before)
	afterMap := ToMap(after)

	plan := RollbackPlan{
		OperationID: operationID,
		CreatedAt:   time.Now().UTC(),
	}

	for k, newVal := range afterMap {
		oldVal, existed := beforeMap[k]
		if !existed || oldVal != newVal {
			plan.Entries = append(plan.Entries, RollbackEntry{
				Key:      k,
				OldValue: oldVal,
				NewValue: newVal,
				HadKey:   existed,
			})
		}
	}

	return plan
}

// Apply reverts the given entries slice to the state captured in the RollbackPlan.
func (p RollbackPlan) Apply(current []Entry) ([]Entry, error) {
	if len(p.Entries) == 0 {
		return current, fmt.Errorf("rollback plan %q has no entries", p.OperationID)
	}

	result := make([]Entry, 0, len(current))
	reverted := make(map[string]bool)

	for _, rb := range p.Entries {
		reverted[rb.Key] = true
	}

	for _, e := range current {
		if !reverted[e.Key] {
			result = append(result, e)
			continue
		}
	}

	for _, rb := range p.Entries {
		if rb.HadKey {
			result = append(result, Entry{Key: rb.Key, Value: rb.OldValue})
		}
		// If HadKey is false, the key is simply dropped (it didn't exist before).
	}

	return result, nil
}
