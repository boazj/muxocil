// Package utils contains common utility functions for use across muxocil
package utils

import (
	"io/fs"
	"os"
	"path/filepath"
)

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
