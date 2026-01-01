// Package sliceutil provides generic utility functions not available in stdlib.
// For most slice operations, use the standard library "slices" package directly.
// This package only contains functions that have no stdlib equivalent.
package sliceutil

// Unique returns a new slice with duplicate elements removed.
// Preserves the order of first occurrence.
// The original slice is not modified.
//
// Note: Unlike slices.Compact, this works on unsorted slices and preserves order.
func Unique[T comparable](slice []T) []T {
	if slice == nil {
		return nil
	}
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// FindMapKey returns the first key in map m whose value equals val.
// Returns the key and true if found, or zero value and false if not found.
func FindMapKey[K comparable, V comparable](m map[K]V, val V) (K, bool) {
	for k, v := range m {
		if v == val {
			return k, true
		}
	}
	var zero K
	return zero, false
}
