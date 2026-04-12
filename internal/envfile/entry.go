package envfile

// Entry represents a single key-value pair parsed from an .env file.
type Entry struct {
	Key     string
	Value   string
	// Comment holds any inline or preceding comment associated with the entry.
	Comment string
}

// ToMap converts a slice of Entry values into a map keyed by Entry.Key.
// If duplicate keys exist, the last value wins.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// FromMap converts a map into a slice of Entry values.
// The order of entries is non-deterministic.
func FromMap(m map[string]string) []Entry {
	entries := make([]Entry, 0, len(m))
	for k, v := range m {
		entries = append(entries, Entry{Key: k, Value: v})
	}
	return entries
}

// Keys returns a slice of all keys from the provided entries.
func Keys(entries []Entry) []string {
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = e.Key
	}
	return keys
}

// Lookup finds the first Entry with the given key and returns it.
// The second return value indicates whether the key was found.
func Lookup(entries []Entry, key string) (Entry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return Entry{}, false
}
