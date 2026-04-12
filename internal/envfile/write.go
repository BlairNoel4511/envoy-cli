package envfile

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Write serialises the given key/value map to a .env file at path.
// Keys are written in sorted order. Values containing spaces or special
// characters are double-quoted.
func Write(path string, env map[string]string) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		if needsQuoting(v) {
			v = fmt.Sprintf("%q", v)
		}
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(v)
		sb.WriteByte('\n')
	}

	//nolint:gosec
	if err := os.WriteFile(path, []byte(sb.String()), 0o644); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// needsQuoting reports whether a value should be wrapped in double quotes
// when written to an env file.
func needsQuoting(v string) bool {
	if v == "" {
		return false
	}
	for _, c := range v {
		if c == ' ' || c == '\t' || c == '#' || c == '"' || c == '\'' || c == '$' {
			return true
		}
	}
	return false
}
