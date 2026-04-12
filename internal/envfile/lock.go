package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LockEntry represents a pinned key-value pair with metadata.
type LockEntry struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	PinnedAt  time.Time `json:"pinned_at"`
	PinnedBy  string    `json:"pinned_by"`
	Comment   string    `json:"comment,omitempty"`
}

// LockFile holds a set of pinned environment variable entries.
type LockFile struct {
	Version int                   `json:"version"`
	Entries map[string]LockEntry  `json:"entries"`
}

// NewLockFile creates an empty LockFile.
func NewLockFile() *LockFile {
	return &LockFile{
		Version: 1,
		Entries: make(map[string]LockEntry),
	}
}

// Pin adds or updates a pinned entry in the lock file.
func (lf *LockFile) Pin(key, value, pinnedBy, comment string) {
	lf.Entries[key] = LockEntry{
		Key:      key,
		Value:    value,
		PinnedAt: time.Now().UTC(),
		PinnedBy: pinnedBy,
		Comment:  comment,
	}
}

// Unpin removes a pinned entry by key. Returns false if key was not pinned.
func (lf *LockFile) Unpin(key string) bool {
	if _, ok := lf.Entries[key]; !ok {
		return false
	}
	delete(lf.Entries, key)
	return true
}

// IsPinned reports whether the given key is pinned.
func (lf *LockFile) IsPinned(key string) bool {
	_, ok := lf.Entries[key]
	return ok
}

// Get returns the LockEntry for key, or an error if not found.
func (lf *LockFile) Get(key string) (LockEntry, error) {
	e, ok := lf.Entries[key]
	if !ok {
		return LockEntry{}, fmt.Errorf("key %q is not pinned", key)
	}
	return e, nil
}

// SaveLockFile writes the LockFile to disk as JSON.
func SaveLockFile(path string, lf *LockFile) error {
	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal lock file: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadLockFile reads a LockFile from disk. Returns an empty LockFile if not found.
func LoadLockFile(path string) (*LockFile, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NewLockFile(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read lock file: %w", err)
	}
	var lf LockFile
	if err := json.Unmarshal(data, &lf); err != nil {
		return nil, fmt.Errorf("parse lock file: %w", err)
	}
	if lf.Entries == nil {
		lf.Entries = make(map[string]LockEntry)
	}
	return &lf, nil
}
