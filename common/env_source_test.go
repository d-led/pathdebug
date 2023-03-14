package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_empty_var_results_in_empty_list(t *testing.T) {
	const varThatShouldBeEmpty = "SHOULD_BE_EMPTY"
	os.Unsetenv(varThatShouldBeEmpty)
	assert.Equal(t, []string{}, NewEnvSource(varThatShouldBeEmpty).Values())
}

const tmpPathVar = "TMP_PATH"
const sep = string(os.PathListSeparator)

func Test_variable_name_is_available(t *testing.T) {
	assert.Equal(t, tmpPathVar, NewEnvSource(tmpPathVar).Source())
}

func Test_repeated_separators_are_removed(t *testing.T) {
	os.Setenv(tmpPathVar, "a"+sep+sep+"b")
	assert.Equal(t, []string{"a", "b"}, NewEnvSource(tmpPathVar).Values())
}

func Test_duplicates_are_preserved(t *testing.T) {
	os.Setenv(tmpPathVar, "a"+sep+"b"+sep+"a")
	assert.Equal(t, []string{"a", "b", "a"}, NewEnvSource(tmpPathVar).Values())
}
