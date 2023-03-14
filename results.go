package main

import (
	"errors"

	"github.com/zyedidia/generic/multimap"
	"github.com/zyedidia/generic/set"
)

type resultRow struct {
	id           int
	path         string
	expandedPath string
	exists       bool
	isDir        bool
	duplicates   []int
}

func calculateResults(fs filesystem, source valueSource) ([]resultRow, error) {
	values := source.values()

	if len(values) == 0 {
		return nil, errors.New(source.source() + " is empty")
	}

	duplicatePredicate := func(a, b int) bool { return true }
	pathsLookup := multimap.NewMapSet[string](duplicatePredicate)
	for i, path := range values {
		pathKey := fs.getAbsolutePath(path)
		pathsLookup.Put(pathKey, i)
	}

	return calculateResultRows(fs, values, pathsLookup), nil
}

func calculateResultRows(fs filesystem, paths []string, pathsLookup multimap.MultiMap[string, int]) []resultRow {
	res := []resultRow{}
	for index, path := range paths {
		pathKey := fs.getAbsolutePath(path)
		dup := getDuplicatesOf(pathsLookup, pathKey, index)
		exists, isdir := fs.pathStatus(pathKey)
		res = append(res, resultRow{
			id:           index + 1,
			path:         path,
			expandedPath: pathKey,
			exists:       exists,
			isDir:        isdir,
			duplicates:   dup,
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
