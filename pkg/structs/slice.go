package structs

func Map[T any, R any](collection []T, f func(T) R) []R {
	if collection == nil {
		return nil
	}

	result := make([]R, len(collection))
	for i, item := range collection {
		result[i] = f(item)
	}

	return result
}

func MapWithError[T any, R any](collection []T, f func(T) (R, error)) ([]R, error) {
	if collection == nil {
		return nil, nil
	}

	result := make([]R, len(collection))
	for i, item := range collection {
		v, err := f(item)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}

	return result, nil
}

func FilterSlice[V any](collection []V, filter func(V) bool) []V {
	if collection == nil {
		return nil
	}

	result := make([]V, 0, len(collection))
	for _, item := range collection {
		if !filter(item) {
			continue
		}
		result = append(result, item)
	}

	return result
}

func ForEach[T any](collection []T, f func(T)) {
	for _, item := range collection {
		f(item)
	}
}

func SliceToMap[K comparable, V any](collection []V, key func(V) K) map[K]V {
	if collection == nil {
		return nil
	}

	result := make(map[K]V, len(collection))
	for _, item := range collection {
		result[key(item)] = item
	}

	return result
}

func SliceToMapOfStruct[K comparable, V any](collection []V, key func(V) K) map[K]struct{} {
	if collection == nil {
		return nil
	}

	result := make(map[K]struct{}, len(collection))
	for _, item := range collection {
		result[key(item)] = struct{}{}
	}

	return result
}
