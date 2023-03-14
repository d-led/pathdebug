package render

import (
	"strings"
	"testing"

	"github.com/d-led/pathdebug/common"
	"github.com/stretchr/testify/assert"
)

func Test_relative_paths_are_expanded(t *testing.T) {
	var b strings.Builder
	RenderTable(&b, []common.ResultRow{
		{Id: 42, Path: "/ok", Duplicates: []int{1, 2}, IsDir: true, Exists: true},
		{Id: 33, Path: "/file", IsDir: false, Exists: true},
		{Id: 33, Path: "/not-ok", IsDir: false},
	})
	table := b.String()

	assert.Contains(t, table, "/ok")
	assert.Contains(t, table, "/file")
	assert.Contains(t, table, "/not-ok")
	assert.Contains(t, table, "42")
	assert.Contains(t, table, "1, 2")
	assert.Contains(t, table, "F")
	assert.Contains(t, table, "X")
}
