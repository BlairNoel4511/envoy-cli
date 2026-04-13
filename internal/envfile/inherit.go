package envfile

// InheritOptions controls how inheritance between profiles is applied.
type InheritOptions struct {
	Overwrite bool // overwrite existing keys in child
	SkipSensitive bool // skip sensitive keys from parent
}

// InheritResult describes the outcome of an inherit operation.
type InheritResult struct {
	Added   []string
	Skipped []string
	Overwritten []string
}

// Inherit merges entries from a parent profile into a child slice of entries.
// Keys already present in child are skipped unless Overwrite is set.
// Sensitive keys from parent are skipped when SkipSensitive is set.
func Inherit(parent, child []Entry, opts InheritOptions) ([]Entry, InheritResult) {
	result := InheritResult{}
	childMap := make(map[string]int) // key -> index in child
	for i, e := range child {
		childMap[e.Key] = i
	}

	output := make([]Entry, len(child))
	copy(output, child)

	for _, pe := range parent {
		if opts.SkipSensitive && IsSensitive(pe.Key) {
			result.Skipped = append(result.Skipped, pe.Key)
			continue
		}

		if idx, exists := childMap[pe.Key]; exists {
			if !opts.Overwrite {
				result.Skipped = append(result.Skipped, pe.Key)
				continue
			}
			if output[idx].Value == pe.Value {
				result.Skipped = append(result.Skipped, pe.Key)
				continue
			}
			output[idx] = pe
			result.Overwritten = append(result.Overwritten, pe.Key)
		} else {
			output = append(output, pe)
			childMap[pe.Key] = len(output) - 1
			result.Added = append(result.Added, pe.Key)
		}
	}

	return output, result
}

// InheritSummary returns a human-readable summary of an InheritResult.
func InheritSummary(r InheritResult) string {
	var parts []string
	if len(r.Added) > 0 {
		parts = append(parts, fmt.Sprintf("%d added", len(r.Added)))
	}
	if len(r.Overwritten) > 0 {
		parts = append(parts, fmt.Sprintf("%d overwritten", len(r.Overwritten)))
	}
	if len(r.Skipped) > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", len(r.Skipped)))
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, ", ")
}
