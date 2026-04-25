package envfile

import "strings"

// Tag represents a label attached to an env entry.
type Tag struct {
	Key   string
	Label string
}

// TagStore holds tags indexed by env key.
type TagStore struct {
	tags map[string][]string
}

// NewTagStore creates an empty TagStore.
func NewTagStore() *TagStore {
	return &TagStore{tags: make(map[string][]string)}
}

// Add attaches a label to the given env key.
func (ts *TagStore) Add(key, label string) {
	key = strings.TrimSpace(key)
	label = strings.TrimSpace(label)
	if key == "" || label == "" {
		return
	}
	for _, existing := range ts.tags[key] {
		if existing == label {
			return
		}
	}
	ts.tags[key] = append(ts.tags[key], label)
}

// Remove detaches a label from the given env key.
func (ts *TagStore) Remove(key, label string) {
	current := ts.tags[key]
	updated := current[:0]
	for _, l := range current {
		if l != label {
			updated = append(updated, l)
		}
	}
	ts.tags[key] = updated
}

// Get returns all labels for a given env key.
func (ts *TagStore) Get(key string) []string {
	return ts.tags[key]
}

// HasTag reports whether a key has a specific label.
func (ts *TagStore) HasTag(key, label string) bool {
	for _, l := range ts.tags[key] {
		if l == label {
			return true
		}
	}
	return false
}

// KeysWithTag returns all env keys that carry the given label.
func (ts *TagStore) KeysWithTag(label string) []string {
	var result []string
	for k, labels := range ts.tags {
		for _, l := range labels {
			if l == label {
				result = append(result, k)
				break
			}
		}
	}
	return result
}

// All returns a flat list of Tag structs.
func (ts *TagStore) All() []Tag {
	var out []Tag
	for k, labels := range ts.tags {
		for _, l := range labels {
			out = append(out, Tag{Key: k, Label: l})
		}
	}
	return out
}

// RemoveKey deletes all labels associated with the given env key.
func (ts *TagStore) RemoveKey(key string) {
	delete(ts.tags, key)
}
