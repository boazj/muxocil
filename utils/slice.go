package utils

func Times[T any](times int, vs ...T) []T {
	s := make([]T, times)
	for range times {
		s = append(s, vs...)
	}
	return s
}

func AppendTimes[T any](slice []T, times int, vs ...T) []T {
	s := slice
	for range times {
		s = append(s, vs...)
	}
	return s
}

func ToStringSlice[T ~string, S ~[]T](slice S) []string {
	result := make([]string, len(slice))
	for i := range slice {
		result[i] = string(slice[i])
	}
	return result
}
