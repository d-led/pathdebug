package common

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

type OsFilesystem struct{}

func (*OsFilesystem) GetAbsolutePath(path string) string {
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
func (f *OsFilesystem) PathStatus(path string) (bool, bool) {
	path = f.GetAbsolutePath(path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, false
	}

	if !fileInfo.IsDir() {
		return true, false
	}

	return true, true
}
