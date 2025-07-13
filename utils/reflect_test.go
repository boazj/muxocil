package utils

import (
	"reflect"
	"strings"
	"testing"
)

func TestAsStringMap(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    map[string]string
		wantErr string
	}{
		{"AsStringMap(string) == err", "name", nil, "trying to convert a non-struct to string map"},
		{"AsStringMap(int) == err", 8, nil, "trying to convert a non-struct to string map"},
		{"AsStringMap(slice) == err", []string{"value"}, nil, "trying to convert a non-struct to string map"},
		{"AsStringMap(map) == err", map[string]string{"v": "value"}, nil, "trying to convert a non-struct to string map"},
		{"AsStringMap(string struct) == string map", struct{ v string }{v: "value"}, map[string]string{"v": "value"}, ""},
		{"AsStringMap(int struct) == empty map", struct{ v int }{v: 8}, map[string]string{}, ""},
		{"AsStringMap(mixed struct) == string map", struct {
			v string
			x int
		}{v: "value", x: 8}, map[string]string{"v": "value"}, ""},
		{"AsStringMap(2 string struct) == 2 string map", struct {
			v string
			x string
		}{v: "value", x: "other"}, map[string]string{"v": "value", "x": "other"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := AsStringMap(tt.input)
			if tt.want != nil && !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
			if tt.wantErr != "" && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("got %v, want %s", err, tt.wantErr)
			}
		})
	}
}
