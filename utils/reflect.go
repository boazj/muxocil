package utils

import (
	"fmt"
	"reflect"
)

func AsStringMap(s any) (map[string]string, error) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("trying to convert a non-struct to string map")
	}
	fields := make(map[string]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if field.Type.Kind() != reflect.String {
			continue
		}
		fields[field.Name] = value.String()
	}
	return fields, nil
}
