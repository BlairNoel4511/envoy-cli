package envfile

import "time"

// DeprecationStatus represents the deprecation state of a key.
type DeprecationStatus string

const (
	DeprecationActive  DeprecationStatus = "active"
	DeprecationWarning DeprecationStatus = "warning"
	DeprecationRemoved DeprecationStatus = "removed"
)

// DeprecatedEntry records deprecation metadata for a single key.
type DeprecatedEntry struct {
	Key        string            `json:"key"`
	Status     DeprecationStatus `json:"status"`
	Reason     string            `json:"reason"`
	ReplacedBy string            `json:"replaced_by,omitempty"`
	Since      time.Time         `json:"since"`
}

// DeprecationStore holds all deprecated key entries.
type DeprecationStore struct {
	Entries map[string]DeprecatedEntry `json:"entries"`
}

// NewDeprecationStore returns an empty DeprecationStore.
func NewDeprecationStore() *DeprecationStore {
	return &DeprecationStore{Entries: make(map[string]DeprecatedEntry)}
}

// Deprecate marks a key as deprecated with the given options.
func (d *DeprecationStore) Deprecate(key, reason, replacedBy string, status DeprecationStatus) {
	d.Entries[key] = DeprecatedEntry{
		Key:        key,
		Status:     status,
		Reason:     reason,
		ReplacedBy: replacedBy,
		Since:      time.Now().UTC(),
	}
}

// IsDeprecated returns true if the key is tracked in the store.
func (d *DeprecationStore) IsDeprecated(key string) bool {
	_, ok := d.Entries[key]
	return ok
}

// Get returns the DeprecatedEntry and whether it exists.
func (d *DeprecationStore) Get(key string) (DeprecatedEntry, bool) {
	e, ok := d.Entries[key]
	return e, ok
}

// Remove removes a key from the deprecation store.
func (d *DeprecationStore) Remove(key string) {
	delete(d.Entries, key)
}

// List returns all deprecated entries as a slice.
func (d *DeprecationStore) List() []DeprecatedEntry {
	out := make([]DeprecatedEntry, 0, len(d.Entries))
	for _, e := range d.Entries {
		out = append(out, e)
	}
	return out
}

// CheckEntries inspects a slice of Entry values and returns any
// DeprecatedEntry records that match keys present in entries.
func (d *DeprecationStore) CheckEntries(entries []Entry) []DeprecatedEntry {
	var hits []DeprecatedEntry
	for _, e := range entries {
		if dep, ok := d.Entries[e.Key]; ok {
			hits = append(hits, dep)
		}
	}
	return hits
}
