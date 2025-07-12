package utils

import (
	"path/filepath"
	"strings"
)

func GetFileNamePart(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func FilenamePart(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
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
		if strings.HasSuffix(path, s) {
			return true
		}
	}
	return false
}
