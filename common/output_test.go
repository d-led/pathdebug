package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var output Output

func Test_bad_assignment_fails(t *testing.T) {
	assert.Error(t, output.Set("bad_output"))
}

func Test_all_ok_assignments(t *testing.T) {
	for _, o := range AllOutputs {
		assert.NoError(t, output.Set(o))
	}
}
