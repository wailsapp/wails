// Package sliceutil provides generic utility functions not available in stdlib.
// For most slice operations, use the standard library "slices" package directly.
// This package only contains functions that have no stdlib equivalent.
package sliceutil

// Unique returns a new slice with duplicate elements removed.
// Preserves the order of first occurrence.
// The original slice is not modified.
//
// Unique returns a new slice containing the first occurrence of each element from the input slice, preserving their original order.
// If the input slice is nil, Unique returns nil.
// The original slice is not modified.
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
// FindMapKey returns the first key in m whose value equals val.
// If no such key exists it returns the zero value of K and false. If multiple keys map to val, the returned key depends on Go's map iteration order.
func FindMapKey[K comparable, V comparable](m map[K]V, val V) (K, bool) {
	for k, v := range m {
		if v == val {
			return k, true
		}
	}
	var zero K
	return zero, false
}