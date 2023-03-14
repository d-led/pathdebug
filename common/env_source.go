package common

import (
	"os"
	"strings"
)

type EnvSource struct {
	name string
}

func NewEnvSource(varName string) *EnvSource {
	return &EnvSource{varName}
}

func (s *EnvSource) Values() []string {
	paths := strings.FieldsFunc(s.Orig(), func(c rune) bool {
		return c == os.PathListSeparator
	})
	return paths
}

func (s *EnvSource) Orig() string {
	return os.Getenv(s.name)
}

func (s *EnvSource) Source() string {
	return s.name
}
