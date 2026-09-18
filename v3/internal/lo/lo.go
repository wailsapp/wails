package lo

// Associate converts a slice to a map by applying keyFn to each element.
func Associate[T any, K comparable, V any](collection []T, keyFn func(T) (K, V)) map[K]V {
	result := make(map[K]V, len(collection))
	for _, item := range collection {
		k, v := keyFn(item)
		result[k] = v
	}
	return result
}

// Contains reports whether v is present in collection.
func Contains[T comparable](collection []T, v T) bool {
	for _, item := range collection {
		if item == v {
			return true
		}
	}
	return false
}

// ContainsBy reports whether any element of collection satisfies predicate.
func ContainsBy[T any](collection []T, predicate func(T) bool) bool {
	for _, item := range collection {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Find returns the first element satisfying predicate and true, or the zero
// value and false if no element matches.
func Find[T any](collection []T, predicate func(T) bool) (T, bool) {
	for _, item := range collection {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// Keys returns the keys of the map in an unspecified order.
func Keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// Must returns val and panics if err is non-nil.
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// Ternary returns ifTrue when condition is true, otherwise ifFalse.
func Ternary[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}

// Without returns a copy of collection with all occurrences of exclude removed.
func Without[T comparable](collection []T, exclude ...T) []T {
	result := make([]T, 0, len(collection))
	for _, item := range collection {
		excluded := false
		for _, ex := range exclude {
			if item == ex {
				excluded = true
				break
			}
		}
		if !excluded {
			result = append(result, item)
		}
	}
	return result
}
