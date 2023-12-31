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
	fs              Filesystem
	source          ValueSource
	pathsLookup     positionLookup
	candidateSource CandidateSource
}

func NewCustomResultsCalculator(fs Filesystem, source ValueSource, candidateSource CandidateSource) *ResultsCalculator {
	return &ResultsCalculator{
		fs:              fs,
		source:          source,
		pathsLookup:     multimap.NewMapSet[string](duplicatePredicate),
		candidateSource: candidateSource,
	}
}

func NewResultsCalculator(fs Filesystem, source ValueSource) *ResultsCalculator {
	return NewCustomResultsCalculator(fs, source, NewNixCandidateSource(fs, source.Source()))
}

func (r *ResultsCalculator) CalculateResults() ([]ResultRow, error) {
	if len(r.source.Values()) == 0 {
		return nil, errors.New(r.source.Source() + " is empty")
	}

	r.calculatePathPositionLookup()

	return r.calculateResultRows(), nil
}

func (r *ResultsCalculator) calculatePathPositionLookup() {
	for i, path := range r.source.Values() {
		pathKey := r.fs.GetAbsolutePath(path)
		r.pathsLookup.Put(pathKey, i)
	}
}

func (r *ResultsCalculator) calculateResultRows() []ResultRow {
	paths := r.source.Values()
	res := []ResultRow{}
	for index, path := range paths {
		pathKey := r.fs.GetAbsolutePath(path)
		dup := getDuplicatesOf(r.pathsLookup, pathKey, index)
		exists, isdir := r.fs.PathStatus(pathKey)
		candidateSources := r.getCandidateSourcesFor(path)
		res = append(res, ResultRow{
			Id:               index + 1,
			Path:             path,
			ExpandedPath:     pathKey,
			Exists:           exists,
			IsDir:            isdir,
			Duplicates:       dup,
			CandidateSources: candidateSources,
		})
	}
	return res
}

func (r *ResultsCalculator) getCandidateSourcesFor(path string) []string {
	res := r.candidateSource.WhereSet(path)
	if res == nil {
		return nil
	}
	candidateSources := []string{}
	for _, source := range res.WhereSet {
		candidateSources = append(candidateSources, source.Original)
	}
	return candidateSources
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
