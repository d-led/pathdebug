package common

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"go.spiff.io/expand"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashset"
)

type NixCandidateSource struct {
	fs                          Filesystem
	sources                     map[string]*PathSetIn
	expandedPathsAlreadyCrawled map[string]bool
	key                         string
}

func NewNixCandidateSource(fs Filesystem, key string) CandidateSource {
	res := &NixCandidateSource{
		fs:                          fs,
		sources:                     map[string]*PathSetIn{},
		expandedPathsAlreadyCrawled: map[string]bool{},
		key:                         key,
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
	ForEachKnownPath(s.crawlSource)
}

func (s *NixCandidateSource) crawlSource(originalSource, expandedSource string) {
	if s.expandedPathsAlreadyCrawled[expandedSource] || isDir(expandedSource) {
		return
	}
	// do not crawl this file again
	s.expandedPathsAlreadyCrawled[expandedSource] = true

	// try getting the contents
	input, err := readAllText(expandedSource)
	if err != nil {
		return
	}

	ForEachVariableAssignment(s.key, input, func(value string) {
		harvestedPaths := s.harvestPaths(value)
		for _, harvestedPath := range harvestedPaths {
			s.tryUpdatePathMap(harvestedPath, originalSource, expandedSource)
		}
	})

	sourcedQueue := []string{}

	ForEachSourcedScript(input, func(sourcedPath string) {
		normalizedSourcedPath := s.fs.GetAbsolutePath(sourcedPath)
		if s.expandedPathsAlreadyCrawled[normalizedSourcedPath] {
			return
		}
		sourcedQueue = append(sourcedQueue, sourcedPath)
	})

	for {
		if len(sourcedQueue) == 0 {
			break
		}
		next := sourcedQueue[0]
		sourcedQueue = sourcedQueue[1:]
		s.crawlSource(next, s.fs.GetAbsolutePath(next))
	}

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
	input = fixHomeExpansion(input)
	input = strings.TrimSpace(input)
	input = strings.Trim(input, fmt.Sprintf("%v\"'", os.PathListSeparator))
	inputExpanded := expand.Expand(input, func(key string) (value string, ok bool) {
		switch key {
		// do not expand the desired variable itself
		case s.key:
			return "", false
		default:
			return os.LookupEnv(key)
		}
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

func isDir(p string) bool {
	f, err := os.Stat(p)
	if err != nil {
		return false
	}
	return f.IsDir()
}
