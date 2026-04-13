package envfile

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a value.
type TransformFunc func(value string) (string, error)

// TransformOp represents a named transformation operation.
type TransformOp string

const (
	TransformUppercase  TransformOp = "uppercase"
	TransformLowercase  TransformOp = "lowercase"
	TransformTrimSpace  TransformOp = "trim"
	TransformTrimPrefix TransformOp = "trim_prefix"
	TransformTrimSuffix TransformOp = "trim_suffix"
	TransformReplace    TransformOp = "replace"
)

// TransformOptions controls how Transform behaves.
type TransformOptions struct {
	Op          TransformOp
	Arg1        string // used for prefix/suffix/replace old
	Arg2        string // used for replace new
	Keys        []string // if empty, applies to all keys
	SkipSensitive bool
}

// TransformResult records the outcome for a single entry.
type TransformResult struct {
	Key      string
	OldValue string
	NewValue string
	Changed  bool
	Skipped  bool
	Reason   string
}

// Transform applies a transformation operation to env entries.
func Transform(entries []Entry, opts TransformOptions) ([]Entry, []TransformResult, error) {
	fn, err := resolveTransformFunc(opts)
	if err != nil {
		return nil, nil, err
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	results := make([]TransformResult, 0, len(entries))
	out := make([]Entry, 0, len(entries))

	for _, e := range entries {
		res := TransformResult{Key: e.Key, OldValue: e.Value}

		if len(keySet) > 0 && !keySet[e.Key] {
			res.Skipped = true
			res.Reason = "not in key list"
			out = append(out, e)
			results = append(results, res)
			continue
		}

		if opts.SkipSensitive && IsSensitive(e.Key) {
			res.Skipped = true
			res.Reason = "sensitive key skipped"
			out = append(out, e)
			results = append(results, res)
			continue
		}

		newVal, err := fn(e.Value)
		if err != nil {
			return nil, nil, fmt.Errorf("transform key %q: %w", e.Key, err)
		}

		res.NewValue = newVal
		res.Changed = newVal != e.Value
		e.Value = newVal
		out = append(out, e)
		results = append(results, res)
	}

	return out, results, nil
}

func resolveTransformFunc(opts TransformOptions) (TransformFunc, error) {
	switch opts.Op {
	case TransformUppercase:
		return func(v string) (string, error) { return strings.ToUpper(v), nil }, nil
	case TransformLowercase:
		return func(v string) (string, error) { return strings.ToLower(v), nil }, nil
	case TransformTrimSpace:
		return func(v string) (string, error) { return strings.TrimSpace(v), nil }, nil
	case TransformTrimPrefix:
		return func(v string) (string, error) { return strings.TrimPrefix(v, opts.Arg1), nil }, nil
	case TransformTrimSuffix:
		return func(v string) (string, error) { return strings.TrimSuffix(v, opts.Arg1), nil }, nil
	case TransformReplace:
		return func(v string) (string, error) { return strings.ReplaceAll(v, opts.Arg1, opts.Arg2), nil }, nil
	default:
		return nil, fmt.Errorf("unknown transform op: %q", opts.Op)
	}
}
