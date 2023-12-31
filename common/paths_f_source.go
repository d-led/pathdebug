package common

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

const pathListsDir = "/etc/paths.d/"

func ForEachPathsDPath(fn func(source string, path string)) {
	_ = filepath.Walk(pathListsDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			input, err := readAllText(path)
			if err != nil {
				log.Print(err)
				return nil
			}

			lines := strings.Split(string(input), "\n")
			for _, line := range lines {
				trimmedPath := strings.TrimSpace(line)
				if trimmedPath == "" {
					continue
				}
				fn(path, trimmedPath)
			}
			return nil
		},
	)
}
