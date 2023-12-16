package utils

func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T

	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}

func Map[T any, R any](slice []T, transform func(T) R) []R {
	var result []R

	for _, item := range slice {
		result = append(result, transform(item))
	}

	return result
}
