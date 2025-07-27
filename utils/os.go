// Package utils contains common utility functions for use across muxocil
package utils

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type GoOs string

const (
	Aix       GoOs = "aix"
	Android   GoOs = "android"
	Darwin    GoOs = "darwin"
	DragonFly GoOs = "dragonfly"
	Freebsd   GoOs = "freebsd"
	Illumos   GoOs = "illumos"
	Ios       GoOs = "ios"
	Js        GoOs = "js"
	Linux     GoOs = "linux"
	Netbsd    GoOs = "netbsd"
	Openbsd   GoOs = "openbsd"
	Plan9     GoOs = "plan9"
	Solaris   GoOs = "solaris"
	Wasip1    GoOs = "wasip1"
	Windows   GoOs = "windows"
	Unknown   GoOs = ""
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
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
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

func CheckIfCmdInPath(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func ExitError(c ExitCode) {
	os.Exit(int(c))
}

func GOOS() GoOs {
	switch runtime.GOOS {
	case "darwin":
		return Darwin
	case "windows":
		return Windows
	case "linux":
		return Linux
	case "aix":
		return Aix
	case "android":
		return Android
	case "dragonfly":
		return DragonFly
	case "freebsd":
		return Freebsd
	case "illumos":
		return Illumos
	case "ios":
		return Ios
	case "js":
		return Js
	case "netbsd":
		return Netbsd
	case "openbsd":
		return Openbsd
	case "plan9":
		return Plan9
	case "solaris":
		return Solaris
	case "wasip1":
		return Wasip1
	default:
		ExitError(ExitUnknownOs)
		return Unknown
	}
}
