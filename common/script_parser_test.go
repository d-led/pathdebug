package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parsing_scrip_assignments(t *testing.T) {
	input := `
if [[ true ]]; then
    export PATH="/path1:$PATH"
else
    export PATH="${PATH}:/path2"
fi
PATH=$PATH:/path3
path=$path:/path4
TEST=$PATH:/path5
`
	var seenValues string
	ForEachVariableAssignment("PATH", input, func(s string) {
		seenValues += s
	})
	assert.Contains(t, seenValues, "/path1")
	assert.Contains(t, seenValues, "/path2")
	assert.Contains(t, seenValues, "/path3")
	assert.NotContains(t, seenValues, "/path4")
	assert.NotContains(t, seenValues, "/path5")
}
