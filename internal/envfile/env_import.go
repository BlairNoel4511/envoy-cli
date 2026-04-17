package envfile

import (
	"fmt"
	"os"
	"strings"
)

// ImportOptions controls how entries are imported from an external source.
type ImportOptions struct {
	Overwrite      bool
	SkipSensitive  bool
	Prefix         string
	DryRun         bool
}

// ImportResult describes what happened to a single key during import.
type ImportResult struct {
	Key      string
	Value    string
	Status   string // added, skipped, overwritten, dry-run
	Sensitive bool
}

// ImportSummary holds aggregate counts from an import operation.
type ImportSummary struct {
	Added      int
	Skipped    int
	Overwritten int
	Total      int
}

// Import merges entries from src into dst according to opts.
// Returns the updated entries, per-key results, and a summary.
func Import(dst []Entry, src []Entry, opts ImportOptions) ([]Entry, []ImportResult, ImportSummary) {
	dstMap := ToMap(dst)
	var results []ImportResult
	summary := ImportSummary{}

	for _, e := range src {
		key := e.Key
		if opts.Prefix != "" {
			if !strings.HasPrefix(key, opts.Prefix) {
				continue
			}
		}

		sensitive := IsSensitive(key)
		if opts.SkipSensitive && sensitive {
			results = append(results, ImportResult{Key: key, Value: e.Value, Status: "skipped", Sensitive: true})
			summary.Skipped++
			summary.Total++
			continue
		}

		_, exists := dstMap[key]
		switch {
		case opts.DryRun:
			results = append(results, ImportResult{Key: key, Value: e.Value, Status: "dry-run", Sensitive: sensitive})
			summary.Total++
		case !exists:
			dst = append(dst, e)
			dstMap[key] = e.Value
			results = append(results, ImportResult{Key: key, Value: e.Value, Status: "added", Sensitive: sensitive})
			summary.Added++
			summary.Total++
		case opts.Overwrite:
			for i, d := range dst {
				if d.Key == key {
					dst[i].Value = e.Value
					break
				}
			}
			results = append(results, ImportResult{Key: key, Value: e.Value, Status: "overwritten", Sensitive: sensitive})
			summary.Overwritten++
			summary.Total++
		default:
			results = append(results, ImportResult{Key: key, Value: e.Value, Status: "skipped", Sensitive: sensitive})
			summary.Skipped++
			summary.Total++
		}
	}
	return dst, results, summary
}

// ImportFromFile parses a .env file and returns its entries.
func ImportFromFile(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("import: cannot open %s: %w", path, err)
	}
	defer f.Close()
	return Parse(f)
}
