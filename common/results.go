package common

import (
	"errors"

	"github.com/zyedidia/generic/multimap"
	"github.com/zyedidia/generic/set"
)

type ResultRow struct {
	Id           int    `json:"id"`
	Path         string `json:"path"`
	ExpandedPath string `json:"expanded_path"`
	Exists       bool   `json:"exists"`
	IsDir        bool   `json:"is_dir"`
	Duplicates   []int  `json:"duplicates"`
}

func CalculateResults(fs Filesystem, source ValueSource) ([]ResultRow, error) {
	values := source.Values()

	if len(values) == 0 {
		return nil, errors.New(source.Source() + " is empty")
	}

	duplicatePredicate := func(a, b int) bool { return true }
	pathsLookup := multimap.NewMapSet[string](duplicatePredicate)
	for i, path := range values {
		pathKey := fs.GetAbsolutePath(path)
		pathsLookup.Put(pathKey, i)
	}

	return calculateResultRows(fs, values, pathsLookup), nil
}

func calculateResultRows(fs Filesystem, paths []string, pathsLookup multimap.MultiMap[string, int]) []ResultRow {
	res := []ResultRow{}
	for index, path := range paths {
		pathKey := fs.GetAbsolutePath(path)
		dup := getDuplicatesOf(pathsLookup, pathKey, index)
		exists, isdir := fs.PathStatus(pathKey)
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
