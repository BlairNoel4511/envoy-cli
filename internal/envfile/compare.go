package envfile

import (
	"fmt"
	"sort"
)

// CompareResult holds the outcome of comparing two sets of entries.
type CompareResult struct {
	OnlyInLeft  []Entry
	OnlyInRight []Entry
	Different   []ComparedPair
	Identical   []Entry
}

// ComparedPair represents a key that exists in both sets but with different values.
type ComparedPair struct {
	Key        string
	LeftValue  string
	RightValue string
}

// Compare performs a full comparison between two slices of entries,
// returning a structured result showing additions, removals, changes, and matches.
func Compare(left, right []Entry) CompareResult {
	leftMap := ToMap(left)
	rightMap := ToMap(right)

	result := CompareResult{}

	for _, e := range left {
		rv, exists := rightMap[e.Key]
		if !exists {
			result.OnlyInLeft = append(result.OnlyInLeft, e)
		} else if rv != e.Value {
			result.Different = append(result.Different, ComparedPair{
				Key:        e.Key,
				LeftValue:  e.Value,
				RightValue: rv,
			})
		} else {
			result.Identical = append(result.Identical, e)
		}
	}

	for _, e := range right {
		if _, exists := leftMap[e.Key]; !exists {
			result.OnlyInRight = append(result.OnlyInRight, e)
		}
	}

	sort.Slice(result.OnlyInLeft, func(i, j int) bool { return result.OnlyInLeft[i].Key < result.OnlyInLeft[j].Key })
	sort.Slice(result.OnlyInRight, func(i, j int) bool { return result.OnlyInRight[i].Key < result.OnlyInRight[j].Key })
	sort.Slice(result.Different, func(i, j int) bool { return result.Different[i].Key < result.Different[j].Key })
	sort.Slice(result.Identical, func(i, j int) bool { return result.Identical[i].Key < result.Identical[j].Key })

	return result
}

// HasDifferences returns true if the comparison found any differences.
func HasDifferences(r CompareResult) bool {
	return len(r.OnlyInLeft) > 0 || len(r.OnlyInRight) > 0 || len(r.Different) > 0
}

// CompareSummary returns a one-line summary of a CompareResult.
func CompareSummary(r CompareResult) string {
	return fmt.Sprintf("%d added, %d removed, %d changed, %d identical",
		len(r.OnlyInRight), len(r.OnlyInLeft), len(r.Different), len(r.Identical))
}
