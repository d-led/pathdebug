package main

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func getAbsolutePath(path string) string {
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

func statusOf(path string) string {
	path = getAbsolutePath(path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "X"
	}

	if !fileInfo.IsDir() {
		return "F"
	}

	return " "
}
