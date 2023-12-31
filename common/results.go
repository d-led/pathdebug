package common

import (
	"errors"

	"github.com/zyedidia/generic/multimap"
	"github.com/zyedidia/generic/set"
)

type ResultRow struct {
	Id               int      `json:"id"`
	Path             string   `json:"path"`
	ExpandedPath     string   `json:"expanded_path"`
	Exists           bool     `json:"exists"`
	IsDir            bool     `json:"is_dir"`
	Duplicates       []int    `json:"duplicates"`
	CandidateSources []string `json:"candidate_sources,omitempty"`
}

type positionLookup multimap.MultiMap[string, int]

func duplicatePredicate(a, b int) bool { return true }

type ResultsCalculator struct {
	fs          Filesystem
	source      ValueSource
	pathsLookup positionLookup
}

func NewResultsCalculator(fs Filesystem, source ValueSource) *ResultsCalculator {
	return &ResultsCalculator{
		fs:          fs,
		source:      source,
		pathsLookup: multimap.NewMapSet[string](duplicatePredicate),
	}
}

func (r *ResultsCalculator) CalculateResults() ([]ResultRow, error) {
	values := r.source.Values()

	if len(values) == 0 {
		return nil, errors.New(r.source.Source() + " is empty")
	}

	for i, path := range values {
		pathKey := r.fs.GetAbsolutePath(path)
		r.pathsLookup.Put(pathKey, i)
	}

	return r.calculateResultRows(), nil
}

func (r *ResultsCalculator) calculateResultRows() []ResultRow {
	paths := r.source.Values()
	res := []ResultRow{}
	for index, path := range paths {
		pathKey := r.fs.GetAbsolutePath(path)
		dup := getDuplicatesOf(r.pathsLookup, pathKey, index)
		exists, isdir := r.fs.PathStatus(pathKey)
		res = append(res, ResultRow{
			Id:           index + 1,
			Path:         path,
			ExpandedPath: pathKey,
			Exists:       exists,
			IsDir:        isdir,
			Duplicates:   dup,
		})
	}
	return res
}

// returns the IDs of the duplicate indices (== index+1)
func getDuplicatesOf(pathsLookup multimap.MultiMap[string, int], pathKey string, index int) []int {
	instances := pathsLookup.Get(pathKey)
	if len(instances) < 2 {
		return []int{}
	}
	s := set.NewMapset(instances...)
	s.Remove(index)
	dup := []int{}
	s.Each(func(key int) {
		dup = append(dup, key+1)
	})
	return dup
}
