package utils

import "github.com/tiendc/gofn"

func MapFilter[K comparable, V any](m map[K]V, filterFunc func(v V) bool) map[K]V {
	e := gofn.MapEntries(m)
	result := make(map[K]V, 0)
	for _, v := range e {
		if filterFunc(v.Elem2) {
			result[v.Elem1] = v.Elem2
		}
	}
	return result
}

func MapFilterValues[K comparable, V any](m map[K]V, filterFunc func(v V) bool) []V {
	e := gofn.MapValues(m)
	result := make([]V, 0)
	for _, v := range e {
		if filterFunc(v) {
			result = append(result, v)
		}
	}
	return result
}
