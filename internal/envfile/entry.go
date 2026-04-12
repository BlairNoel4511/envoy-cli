package envfile

// Entry represents a single key-value pair from an .env file.
type Entry struct {
	Key   string
	Value string
}

// ToMap converts a slice of Entry to a map[string]string.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// FromMap converts a map[string]string to a slice of Entry.
// Order is non-deterministic; callers should sort if needed.
func FromMap(m map[string]string) []Entry {
	entries := make([]Entry, 0, len(m))
	for k, v := range m {
		entries = append(entries, Entry{Key: k, Value: v})
	}
	return entries
}

// Keys returns the keys of all entries in order.
func Keys(entries []Entry) []string {
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = e.Key
	}
	return keys
}

// Lookup returns the value for the given key and whether it was found.
func Lookup(entries []Entry, key string) (string, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}
