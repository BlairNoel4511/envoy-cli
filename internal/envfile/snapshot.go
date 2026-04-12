package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of an env file at a point in time.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Entries   []Entry           `json:"entries"`
	Checksum  string            `json:"checksum"`
}

// Entry represents a single key-value pair from an env file.
type Entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewSnapshot creates a Snapshot from a parsed env map.
func NewSnapshot(source string, env map[string]string) Snapshot {
	entries := make([]Entry, 0, len(env))
	for k, v := range env {
		entries = append(entries, Entry{Key: k, Value: v})
	}
	return Snapshot{
		Timestamp: time.Now().UTC(),
		Source:    source,
		Entries:   entries,
		Checksum:  checksumEntries(entries),
	}
}

// SaveSnapshot writes a Snapshot to a JSON file at the given path.
func SaveSnapshot(path string, snap Snapshot) error {
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}
	return nil
}

// LoadSnapshot reads and parses a Snapshot from a JSON file.
func LoadSnapshot(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read failed: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: parse failed: %w", err)
	}
	return snap, nil
}

// ToMap converts a Snapshot's entries back into a key-value map.
func (s Snapshot) ToMap() map[string]string {
	m := make(map[string]string, len(s.Entries))
	for _, e := range s.Entries {
		m[e.Key] = e.Value
	}
	return m
}

// checksumEntries produces a simple deterministic hash of entries.
func checksumEntries(entries []Entry) string {
	h := uint32(2166136261)
	for _, e := range entries {
		for _, c := range e.Key + "=" + e.Value + "\n" {
			h ^= uint32(c)
			h *= 16777619
		}
	}
	return fmt.Sprintf("%08x", h)
}
