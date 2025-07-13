package utils

import (
	"path/filepath"
	"strings"
)

func GetFileNamePart(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func IsYaml(path string) bool {
	return HasSuffix(path, ".yaml", ".yml")
}

func IsToml(path string) bool {
	return HasSuffix(path, ".toml")
}

func HasSuffix(path string, suffixs ...string) bool {
	if len(suffixs) == 0 {
		panic("expects at least one suffix")
	}
	for _, s := range suffixs {
		if s != "" && strings.HasSuffix(path, s) {
			return true
		}
	}
	return false
}
