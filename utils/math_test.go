package utils

import (
	"testing"
)

func TestAbsInt(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"AbsInt(10) == 10", 10, 10},
		{"AbsInt(1) == 1", 1, 1},
		{"AbsInt(0) == 0", 0, 0},
		{"AbsInt(-1) == 1", -1, 1},
		{"AbsInt(-10) == 10", -10, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := AbsInt(tt.input)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

func TestOrdinal(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  string
	}{
		{"Ordinal(10) == \"10th\"", 10, "10th"},
		{"Ordinal(1) == \"1st\"", 1, "1st"},
		{"Ordinal(11) == \"11th\"", 11, "11th"},
		{"Ordinal(12) == \"12th\"", 12, "12th"},
		{"Ordinal(13) == \"13th\"", 13, "13th"},
		{"Ordinal(2) == \"2nd\"", 2, "2nd"},
		{"Ordinal(3) == \"3rd\"", 3, "3rd"},
		{"Ordinal(9) == \"9th\"", 9, "9th"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := Ordinal(tt.input)
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}
