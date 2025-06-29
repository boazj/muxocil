package utils

import (
	"os"
)

func GetEnvOr(primary string, secondary string) string {
	term, ok := os.LookupEnv(primary)
	if !ok {
		term = os.Getenv(secondary)
	}
	return term
}

func IsEnvExists(name string) bool {
	_, exists := os.LookupEnv(name)
	return exists
}
