package common

import (
	"strings"
)

// CompareStringSlices compares to slices of strings. It
// returns true if they are the same and false if they are
// different. It takes two slices of strings as input.
func CompareStringSlices(a []string, b []string) bool {
	lenA := len(a)
	lenB := len(b)

	// Make sure that the number of strings is the same
	if lenA != lenB {
		return false
	}

	// For each string
	for i := range a {
		// Make sure they are the same
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// CompareStringSlicesRandom compares to slices of strings in which
// the contents of a string will be in random order. It
// returns true if they are the same and false if they are
// different. It takes a slice of strings (with the randomly-ordered
// contents) and a 2d slice of strings (with the contents broken into
// pieces).
func CompareStringSlicesRandom(a []string, b [][]string) bool {
	lenA := len(a)
	lenB := len(b)

	// Make sure that the number of strings is the same
	if lenA != lenB {
		return false
	}

	// For each string
	for i := range a {
		currentBLen := 0

		// For each piece of the string
		for j := range b[i] {
			// Make sure the piece is in the string
			if !strings.Contains(a[i], b[i][j]) {
				return false
			}

			currentBLen = currentBLen + len(b[i][j])
		}

		// Make sure the length of the string is the same
		// as the length of the pieces
		if len(a[i]) != currentBLen {
			return false
		}
	}

	return true
}
