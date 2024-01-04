package common

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

var knownScriptDirs = []string{
	"/etc/profile.d/",
}

var allowedExtensions = []string{
	".sh",
}

func ForEachScriptsDPath(fn func(originalSource, expandedSource string)) {
	for _, knownScriptDir := range knownScriptDirs {
		_ = filepath.Walk(knownScriptDir,
			func(originalSource string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() || fileNotAllowed(originalSource) {
					return nil
				}

				expanded, err := homedir.Expand(originalSource)
				if err != nil {
					return nil
				}

				fn(originalSource, expanded)

				return nil
			},
		)
	}
}

func fileNotAllowed(p string) bool {
	return !extensionAllowed(p)
}

func extensionAllowed(p string) bool {
	e := filepath.Ext(p)
	for _, allowedExt := range allowedExtensions {
		if allowedExt == e {
			return true
		}
	}
	return false
}
