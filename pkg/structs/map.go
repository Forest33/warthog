// Package structs provides basic operations with maps and slices
package structs

// Keys returns a slice with map keys
func Keys[K comparable, V any](in map[K]V) []K {
	if in == nil {
		return nil
	}

	result := make([]K, 0, len(in))
	for k := range in {
		result = append(result, k)
	}

	return result
}

// FilterMap filters the map
func FilterMap[K comparable, V any](collection map[K]V, filter func(V) bool) map[K]V {
	if collection == nil {
		return nil
	}

	result := make(map[K]V, len(collection))
	for k, v := range collection {
		if !filter(v) {
			continue
		}
		result[k] = v
	}

	return result
}
