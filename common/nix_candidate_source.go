package common

import "runtime"

type NixCandidateSource struct {
	fs      Filesystem
	sources map[string]*PathSetIn
}

func NewNixCandidateSource(fs Filesystem) CandidateSource {
	return &NixCandidateSource{
		fs,
		map[string]*PathSetIn{},
	}
}

func (s *NixCandidateSource) WhereSet(somePath string) *PathSetIn {
	if runtime.GOOS == "windows" {
		return nil
	}
	normalizedPath := s.fs.GetAbsolutePath(somePath)
	return s.sources[normalizedPath]
}
