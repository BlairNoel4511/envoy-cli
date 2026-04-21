package envfile

// CommentOptions controls how comments are applied to entries.
type CommentOptions struct {
	Overwrite bool
}

// CommentResult holds the result of a comment operation on a single entry.
type CommentResult struct {
	Key     string
	Comment string
	Status  string // "added", "updated", "unchanged", "not_found"
}

// CommentSummary holds aggregate counts from a Comment operation.
type CommentSummary struct {
	Added     int
	Updated   int
	Unchanged int
	NotFound  int
}

// Comment sets or updates an inline comment for the given key.
func Comment(entries []Entry, key, comment string, opts CommentOptions) ([]Entry, CommentResult) {
	for i, e := range entries {
		if e.Key != key {
			continue
		}
		if e.Comment == comment {
			return entries, CommentResult{Key: key, Comment: comment, Status: "unchanged"}
		}
		if e.Comment != "" && !opts.Overwrite {
			return entries, CommentResult{Key: key, Comment: e.Comment, Status: "unchanged"}
		}
		status := "added"
		if e.Comment != "" {
			status = "updated"
		}
		entries[i].Comment = comment
		return entries, CommentResult{Key: key, Comment: comment, Status: status}
	}
	return entries, CommentResult{Key: key, Comment: "", Status: "not_found"}
}

// CommentMany applies comments to multiple keys using a map of key -> comment.
func CommentMany(entries []Entry, comments map[string]string, opts CommentOptions) ([]Entry, []CommentResult, CommentSummary) {
	var results []CommentResult
	var sum CommentSummary
	for key, comment := range comments {
		var res CommentResult
		entries, res = Comment(entries, key, comment, opts)
		results = append(results, res)
		switch res.Status {
		case "added":
			sum.Added++
		case "updated":
			sum.Updated++
		case "unchanged":
			sum.Unchanged++
		case "not_found":
			sum.NotFound++
		}
	}
	return entries, results, sum
}
