package utils

import (
	"fmt"
	"math"
)

//nolint:mnd
func Ordinal(number int) string {
	ord := "th"
	switch AbsInt(number) % 100 {
	case 11, 12, 13:
		ord = "th"
	default:
		switch AbsInt(number) % 10 {
		case 1:
			ord = "st"
		case 2:
			ord = "nd"
		case 3:
			ord = "rd"
		}
	}
	return fmt.Sprintf("%d%s", number, ord)
}

func AbsInt(number int) int {
	return int(math.Abs(float64(number)))
}
