// Package utils contains common utility functions for use across muxocil
package utils

import (
	"io/fs"
	"os"
	"path/filepath"
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

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, Wrap(err, "failed to get stats for path")
	}

	return fileInfo.IsDir(), err
}

func GetFilesRecursively(root string, filter func(string) bool) []string {
	files := make([]string, 0)
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filter(path) {
			files = append(files, path)
		}
		return nil
	})
	return files
}
