package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_all_known_paths_can_be_expanded(t *testing.T) {
	count := 0
	ForEachKnownPath(func(_, _ string) {
		count++
	})
	assert.Greater(t, count, 0)
	assert.Equal(t, count, len(knownPaths))
}
