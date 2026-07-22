package helper

import "strconv"

// AtoiDefault parses v as a positive integer, returning fallback if v is
// empty, not a valid integer, or parses to a value <= 0.
func AtoiDefault(v string, fallback int) int {
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
