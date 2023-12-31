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

func Test_recursively_walking_sourced_scripts(t *testing.T) {
	input := `
	[ -s "/test1.sh" ] && \. "/test1.sh"
	[ -f /test2 ] && source /test2
	if [ "${BASH-no}" != "no" ]; then
		[ -r ~/.test3 ] && . '~/.test3'
	fi
	./test4
	`

	seenScripts := []string{}
	ForEachSourcedScript(input, func(s string) {
		seenScripts = append(seenScripts, s)
	})

	assert.Contains(t, seenScripts, "/test1.sh")
	assert.Contains(t, seenScripts, "/test2")
	assert.Contains(t, seenScripts, "~/.test3")
	assert.NotContains(t, seenScripts, "./test4")
}
