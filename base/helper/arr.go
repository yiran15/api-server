package helper

import "slices"

func InArray[T comparable](arr []T, val T) bool {
	return slices.Contains(arr, val)
}
