package envfile

import "fmt"

// PatchOp represents a single patch operation type.
type PatchOp string

const (
	PatchSet    PatchOp = "set"
	PatchDelete PatchOp = "delete"
	PatchRename PatchOp = "rename"
)

// PatchInstruction describes a single mutation to apply to an env entry set.
type PatchInstruction struct {
	Op     PatchOp
	Key    string
	Value  string // used by set
	NewKey string // used by rename
}

// PatchResult records the outcome of a single patch instruction.
type PatchResult struct {
	Instruction PatchInstruction
	Applied     bool
	Skipped     bool
	Reason      string
}

// PatchOptions controls how patch operations behave.
type PatchOptions struct {
	Overwrite bool
	DryRun    bool
}

// Patch applies a list of PatchInstructions to entries and returns updated
// entries along with per-instruction results.
func Patch(entries []Entry, instructions []PatchInstruction, opts PatchOptions) ([]Entry, []PatchResult) {
	results := make([]PatchResult, 0, len(instructions))
	working := make([]Entry, len(entries))
	copy(working, entries)

	for _, ins := range instructions {
		result := PatchResult{Instruction: ins}

		switch ins.Op {
		case PatchSet:
			working, result = applySet(working, ins, opts)
		case PatchDelete:
			working, result = applyDelete(working, ins, opts)
		case PatchRename:
			working, result = applyRename(working, ins, opts)
		default:
			result.Skipped = true
			result.Reason = fmt.Sprintf("unknown op: %s", ins.Op)
		}

		results = append(results, result)
	}

	return working, results
}

func applySet(entries []Entry, ins PatchInstruction, opts PatchOptions) ([]Entry, PatchResult) {
	res := PatchResult{Instruction: ins}
	for i, e := range entries {
		if e.Key == ins.Key {
			if e.Value == ins.Value {
				res.Skipped = true
				res.Reason = "identical value"
				return entries, res
			}
			if !opts.Overwrite {
				res.Skipped = true
				res.Reason = "key exists, overwrite disabled"
				return entries, res
			}
			if !opts.DryRun {
				entries[i].Value = ins.Value
			}
			res.Applied = true
			return entries, res
		}
	}
	if !opts.DryRun {
		entries = append(entries, Entry{Key: ins.Key, Value: ins.Value})
	}
	res.Applied = true
	return entries, res
}

func applyDelete(entries []Entry, ins PatchInstruction, opts PatchOptions) ([]Entry, PatchResult) {
	res := PatchResult{Instruction: ins}
	for i, e := range entries {
		if e.Key == ins.Key {
			if !opts.DryRun {
				entries = append(entries[:i], entries[i+1:]...)
			}
			res.Applied = true
			return entries, res
		}
	}
	res.Skipped = true
	res.Reason = "key not found"
	return entries, res
}

func applyRename(entries []Entry, ins PatchInstruction, opts PatchOptions) ([]Entry, PatchResult) {
	res := PatchResult{Instruction: ins}
	_, existsNew := Lookup(entries, ins.NewKey)
	if existsNew && !opts.Overwrite {
		res.Skipped = true
		res.Reason = "target key exists, overwrite disabled"
		return entries, res
	}
	for i, e := range entries {
		if e.Key == ins.Key {
			if !opts.DryRun {
				entries[i].Key = ins.NewKey
			}
			res.Applied = true
			return entries, res
		}
	}
	res.Skipped = true
	res.Reason = "key not found"
	return entries, res
}
