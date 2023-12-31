package common

type NixCandidateSource struct {
	fs Filesystem
}

func NewNixCandidateSource(fs Filesystem) CandidateSource {
	return &NixCandidateSource{fs}
}

func (s *NixCandidateSource) WhereSet(somePath string) *PathSetIn {
	return nil
}
