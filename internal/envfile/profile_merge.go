package envfile

// ProfileMergeOptions controls how two profiles are merged.
type ProfileMergeOptions struct {
	// Overwrite replaces existing keys in the base profile with values from the overlay.
	Overwrite bool
	// SkipSensitive prevents sensitive keys from being overwritten even when Overwrite is true.
	SkipSensitive bool
}

// MergeProfiles merges the overlay profile into the base profile and returns a new Profile.
// The returned profile inherits the name and tags of the base.
func MergeProfiles(base, overlay *Profile, opts ProfileMergeOptions) *Profile {
	result := &Profile{
		Name: base.Name,
		Tags: append([]string(nil), base.Tags...),
	}

	existing := make(map[string]int) // key -> index in result.Entries
	for _, e := range base.Entries {
		result.Entries = append(result.Entries, e)
		existing[e.Key] = len(result.Entries) - 1
	}

	for _, e := range overlay.Entries {
		idx, found := existing[e.Key]
		if !found {
			result.Entries = append(result.Entries, e)
			existing[e.Key] = len(result.Entries) - 1
			continue
		}
		if opts.Overwrite {
			if opts.SkipSensitive && IsSensitive(e.Key) {
				continue
			}
			result.Entries[idx] = e
		}
	}

	return result
}

// ProfileToEntries returns the entries of a profile by name from the store.
// Returns nil if the profile does not exist.
func ProfileToEntries(ps *ProfileStore, name string) []Entry {
	p, ok := ps.Get(name)
	if !ok {
		return nil
	}
	return p.Entries
}

// ProfileFromEntries creates a Profile from a slice of entries with the given name.
func ProfileFromEntries(name string, entries []Entry) *Profile {
	return &Profile{
		Name:    name,
		Entries: append([]Entry(nil), entries...),
	}
}
