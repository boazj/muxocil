package utils

func AppendTimes[T any](slice []T, times int, vs ...T) []T {
	s := slice
	for range times {
		s = append(s, vs...)
	}
	return s
}
