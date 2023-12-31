package common

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"go.spiff.io/expand"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashset"
	"mvdan.cc/sh/v3/syntax"
)

type NixCandidateSource struct {
	fs      Filesystem
	sources map[string]*PathSetIn
	key     string
}

func NewNixCandidateSource(fs Filesystem, key string) CandidateSource {
	res := &NixCandidateSource{
		fs:      fs,
		sources: map[string]*PathSetIn{},
		key:     key,
	}
	res.crawlKnownPaths()
	res.crawlPathLists()
	return res
}

func (s *NixCandidateSource) WhereSet(somePath string) *PathSetIn {
	if runtime.GOOS == "windows" {
		return nil
	}
	normalizedPath := s.fs.GetAbsolutePath(somePath)
	return s.sources[normalizedPath]
}

func (s *NixCandidateSource) crawlKnownPaths() {
	ForEachKnownPath(func(originalSource, expandedSource string) {
		// try getting the contents
		input, err := readAllText(expandedSource)
		if err != nil {
			return
		}

		// parse
		r := strings.NewReader(input)
		f, err := syntax.NewParser().Parse(r, "")
		if err != nil {
			return
		}

		syntax.Walk(f, func(node syntax.Node) bool {
			switch x := node.(type) {
			case *syntax.Assign:
				if s.key == x.Name.Value {
					harvestedPaths := s.harvestPaths(input[x.Value.Pos().Offset():x.Value.End().Offset()])
					for _, harvestedPath := range harvestedPaths {
						s.tryUpdatePathMap(harvestedPath, originalSource, expandedSource)
					}
				}
			}
			return true
		})
	})
}

func (s *NixCandidateSource) tryUpdatePathMap(harvestedPath string, originalSource string, expandedSource string) {
	foundNormalizedPath := s.fs.GetAbsolutePath(harvestedPath)
	sourcesForPath := s.sources[foundNormalizedPath]
	if sourcesForPath == nil {
		sourcesForPath = &PathSetIn{
			What:     Location{harvestedPath, foundNormalizedPath},
			WhereSet: []Location{},
		}
	}
	sourcesForPath.WhereSet = appendIfNotInSlice(
		sourcesForPath.WhereSet,
		Location{
			originalSource,
			expandedSource,
		}, func(a, b Location) bool {
			return a.Expanded == b.Expanded && a.Original == b.Original
		})
	s.sources[foundNormalizedPath] = sourcesForPath
}

func (s *NixCandidateSource) crawlPathLists() {
	if runtime.GOOS == "windows" {
		return
	}
	ForEachPathsDPath(func(source, path string) {
		s.tryUpdatePathMap(path, source, source)
	})
}

// input is some path definition
func (s *NixCandidateSource) harvestPaths(input string) []string {
	input = strings.TrimSpace(input)
	input = strings.Trim(input, fmt.Sprintf("%v\"'", os.PathListSeparator))
	inputExpanded := expand.Expand(input, func(key string) (value string, ok bool) {
		// do not expand the desired variable itself
		if key == s.key {
			return "", false
		}
		return os.LookupEnv(key)
	})
	all := strings.Split(inputExpanded, fmt.Sprintf("%c", os.PathListSeparator))
	res := hashset.New[string](uint64(len(all)), g.Equals[string], g.HashString)
	for _, p := range all {
		if p != "" {
			res.Put(p)
		}
	}
	r := []string{}
	res.Each(func(key string) {
		r = append(r, key)
	})
	return r
}

func readAllText(filename string) (string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func appendIfNotInSlice[T any](input []T, value T, eq func(T, T) bool) []T {
	res := input
	found := false

	for _, v := range res {
		if eq(v, value) {
			found = true
			break
		}
	}

	if !found {
		res = append(res, value)
	}

	return res
}
