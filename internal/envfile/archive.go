package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ArchiveEntry represents a versioned snapshot of env entries.
type ArchiveEntry struct {
	ID        string   `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Label     string   `json:"label"`
	Entries   []Entry  `json:"entries"`
}

// Archive holds a collection of archived env snapshots.
type Archive struct {
	Records []ArchiveEntry `json:"records"`
}

// NewArchive creates an empty Archive.
func NewArchive() *Archive {
	return &Archive{}
}

// Add appends a new versioned snapshot to the archive.
func (a *Archive) Add(label string, entries []Entry) ArchiveEntry {
	record := ArchiveEntry{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Entries:   append([]Entry{}, entries...),
	}
	a.Records = append(a.Records, record)
	return record
}

// Get returns an ArchiveEntry by ID, and a bool indicating if it was found.
func (a *Archive) Get(id string) (ArchiveEntry, bool) {
	for _, r := range a.Records {
		if r.ID == id {
			return r, true
		}
	}
	return ArchiveEntry{}, false
}

// List returns all archive records.
func (a *Archive) List() []ArchiveEntry {
	return a.Records
}

// SaveArchive writes the archive to a JSON file.
func SaveArchive(path string, a *Archive) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadArchive reads an archive from a JSON file.
// Returns an empty archive if the file does not exist.
func LoadArchive(path string) (*Archive, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NewArchive(), nil
	}
	if err != nil {
		return nil, err
	}
	var a Archive
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, err
	}
	return &a, nil
}
