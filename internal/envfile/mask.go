package envfile

import (
	"fmt"
	"strings"
)

// MaskOptions controls how values are masked.
type MaskOptions struct {
	// VisibleChars is the number of characters to reveal at the start/end.
	VisibleChars int
	// MaskChar is the character used for masking (default '*').
	MaskChar rune
	// MaskLength is the fixed length of the mask segment (0 = proportional).
	MaskLength int
	// SensitiveOnly restricts masking to sensitive keys only.
	SensitiveOnly bool
}

// DefaultMaskOptions returns sensible defaults for masking.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		VisibleChars:  2,
		MaskChar:      '*',
		MaskLength:    6,
		SensitiveOnly: true,
	}
}

// MaskResult holds the result of masking a single entry.
type MaskResult struct {
	Key      string
	Original string
	Masked   string
	WasMasked bool
}

// MaskValue masks a single string value according to the given options.
func MaskValue(value string, opts MaskOptions) string {
	if value == "" {
		return ""
	}
	visible := opts.VisibleChars
	if visible < 0 {
		visible = 0
	}
	maskLen := opts.MaskLength
	if maskLen <= 0 {
		maskLen = max(4, len(value)-visible*2)
	}
	maskChar := opts.MaskChar
	if maskChar == 0 {
		maskChar = '*'
	}
	if len(value) <= visible*2 {
		return strings.Repeat(string(maskChar), maskLen)
	}
	prefix := value[:visible]
	suffix := value[len(value)-visible:]
	mask := strings.Repeat(string(maskChar), maskLen)
	return fmt.Sprintf("%s%s%s", prefix, mask, suffix)
}

// MaskEntries applies masking to a slice of entries, returning MaskResult per entry.
func MaskEntries(entries []Entry, opts MaskOptions) []MaskResult {
	results := make([]MaskResult, 0, len(entries))
	for _, e := range entries {
		sensitive := IsSensitive(e.Key)
		if opts.SensitiveOnly && !sensitive {
			results = append(results, MaskResult{
				Key:      e.Key,
				Original: e.Value,
				Masked:   e.Value,
				WasMasked: false,
			})
			continue
		}
		masked := MaskValue(e.Value, opts)
		results = append(results, MaskResult{
			Key:      e.Key,
			Original: e.Value,
			Masked:   masked,
			WasMasked: masked != e.Value,
		})
	}
	return results
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
