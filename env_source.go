package main

import (
	"os"
	"strings"
)

type envSource struct {
	name string
}

func newEnvSource() *envSource {
	if len(os.Args) != 2 {
		failWith("please provide the name of the environment variable to debug")
	}
	return &envSource{os.Args[1]}
}

func (s *envSource) values() []string {
	paths := strings.Split(s.orig(), string(os.PathListSeparator))
	if len(paths) == 1 && paths[0] == "" {
		return []string{}
	}
	return paths
}

func (s *envSource) orig() string {
	return os.Getenv(s.name)
}

func (s *envSource) source() string {
	return s.name
}
