package envfile

// MergeStrategy defines how conflicts are resolved during a merge.
type MergeStrategy int

const (
	// MergePreferLocal keeps local values on conflict.
	MergePreferLocal MergeStrategy = iota
	// MergePreferRemote overwrites local values with remote on conflict.
	MergePreferRemote
	// MergeAddOnly only adds keys missing in local; never overwrites.
	MergeAddOnly
)

// MergeResult holds the merged env map and metadata about the operation.
type MergeResult struct {
	Merged    map[string]string
	Conflicts []string // keys that had conflicting values
	Added     []string // keys added from remote
}

// Merge combines local and remote env maps according to the given strategy.
func Merge(local, remote map[string]string, strategy MergeStrategy) MergeResult {
	result := MergeResult{
		Merged: make(map[string]string),
	}

	// Copy local into merged as the base.
	for k, v := range local {
		result.Merged[k] = v
	}

	for k, remoteVal := range remote {
		localVal, exists := local[k]
		if !exists {
			result.Merged[k] = remoteVal
			result.Added = append(result.Added, k)
			continue
		}

		if localVal == remoteVal {
			continue
		}

		// Conflict: values differ.
		result.Conflicts = append(result.Conflicts, k)
		switch strategy {
		case MergePreferRemote:
			result.Merged[k] = remoteVal
		case MergePreferLocal, MergeAddOnly:
			result.Merged[k] = localVal
		}
	}

	return result
}
