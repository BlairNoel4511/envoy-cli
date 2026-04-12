package envfile

import (
	"fmt"
	"os"
)

// SyncResult holds the outcome of a sync operation.
type SyncResult struct {
	Written  []string
	Skipped  []string
	Overwritten []string
}

// SyncOptions controls how a sync operation behaves.
type SyncOptions struct {
	// Overwrite existing keys in the destination file.
	Overwrite bool
	// DryRun reports what would change without writing.
	DryRun bool
}

// Sync merges src entries into the destination file at destPath.
// If the destination file does not exist it is created.
// It returns a SyncResult describing what changed.
func Sync(destPath string, src map[string]string, opts SyncOptions) (SyncResult, error) {
	var result SyncResult

	dest := make(map[string]string)
	if _, err := os.Stat(destPath); err == nil {
		parsed, err := Parse(destPath)
		if err != nil {
			return result, fmt.Errorf("sync: reading destination: %w", err)
		}
		for _, e := range parsed {
			dest[e.Key] = e.Value
		}
	}

	for k, v := range src {
		if existing, exists := dest[k]; exists {
			if existing == v {
				result.Skipped = append(result.Skipped, k)
				continue
			}
			if !opts.Overwrite {
				result.Skipped = append(result.Skipped, k)
				continue
			}
			dest[k] = v
			result.Overwritten = append(result.Overwritten, k)
		} else {
			dest[k] = v
			result.Written = append(result.Written, k)
		}
	}

	if opts.DryRun {
		return result, nil
	}

	if err := Write(destPath, dest); err != nil {
		return result, fmt.Errorf("sync: writing destination: %w", err)
	}

	return result, nil
}
