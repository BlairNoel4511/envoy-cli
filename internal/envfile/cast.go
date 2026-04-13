package envfile

import (
	"fmt"
	"strconv"
	"strings"
)

// CastResult holds the result of casting an entry's value to a target type.
type CastResult struct {
	Key      string
	Original string
	Casted   string
	Type     string
	Changed  bool
	Error    string
}

// CastOptions controls how values are cast.
type CastOptions struct {
	// Types maps key names to target type: "int", "bool", "float", "string"
	Types map[string]string
	// SkipErrors continues on cast failure instead of returning an error
	SkipErrors bool
}

// Cast attempts to normalize entry values to the specified target types.
// It returns the updated entries and a slice of CastResult for reporting.
func Cast(entries []Entry, opts CastOptions) ([]Entry, []CastResult, error) {
	results := make([]CastResult, 0, len(entries))
	out := make([]Entry, len(entries))
	copy(out, entries)

	for i, e := range out {
		targetType, ok := opts.Types[e.Key]
		if !ok {
			continue
		}

		casted, err := castValue(e.Value, targetType)
		r := CastResult{
			Key:      e.Key,
			Original: e.Value,
			Type:     targetType,
		}

		if err != nil {
			r.Error = err.Error()
			r.Casted = e.Value
			results = append(results, r)
			if !opts.SkipErrors {
				return nil, results, fmt.Errorf("cast failed for key %q: %w", e.Key, err)
			}
			continue
		}

		r.Casted = casted
		r.Changed = casted != e.Value
		results = append(results, r)
		out[i].Value = casted
	}

	return out, results, nil
}

func castValue(val, targetType string) (string, error) {
	switch strings.ToLower(targetType) {
	case "int":
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return val, fmt.Errorf("cannot cast %q to int", val)
		}
		return strconv.Itoa(int(v)), nil
	case "float":
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return val, fmt.Errorf("cannot cast %q to float", val)
		}
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case "bool":
		v, err := strconv.ParseBool(val)
		if err != nil {
			return val, fmt.Errorf("cannot cast %q to bool", val)
		}
		return strconv.FormatBool(v), nil
	case "string":
		return val, nil
	default:
		return val, fmt.Errorf("unknown target type %q", targetType)
	}
}
