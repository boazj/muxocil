package utils

import "fmt"

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return WrapF(err, msg)
}

func WrapF(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}
	arr := make([]any, len(a)+1)
	arr = append(arr, err)
	arr = append(arr, a...)
	return fmt.Errorf(format, arr...)
}
