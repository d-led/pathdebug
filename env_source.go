package main

import (
	"os"
	"strings"
)

type envSource struct {
	name string
}

func newEnvSource(varName string) *envSource {
	return &envSource{varName}
}

func (s *envSource) values() []string {
	paths := strings.FieldsFunc(s.orig(), func(c rune) bool {
		return c == os.PathListSeparator
	})
	return paths
}

func (s *envSource) orig() string {
	return os.Getenv(s.name)
}

func (s *envSource) source() string {
	return s.name
}
