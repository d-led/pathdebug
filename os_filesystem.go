package main

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

type osFilesystem struct{}

func (*osFilesystem) getAbsolutePath(path string) string {
	// homedir is assumed to work correctly
	expandedPath, err := homedir.Expand(path)
	if err == nil {
		path = expandedPath
	}
	absPath, err := filepath.Abs(path)
	if err == nil {
		path = absPath
	}
	return os.ExpandEnv(path)
}

// returns (exists, isDir)
func (f *osFilesystem) pathStatus(path string) (bool, bool) {
	path = f.getAbsolutePath(path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, false
	}

	if !fileInfo.IsDir() {
		return true, false
	}

	return true, true
}
