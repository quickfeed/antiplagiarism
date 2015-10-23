package common

// CompareStringSlices compares to slices of strings. It
// returns true if they are the same and false if they are
// different. It takes two slices of strings as input.
func CompareStringSlices(a []string, b []string) bool {
	lenA := len(a)
	lenB := len(b)

	if lenA != lenB {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
