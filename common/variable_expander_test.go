package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_expanding_with_custom_callback(t *testing.T) {
	input := "${VAR1}/$VAR2/bin"

	res := CustomExpandVariables(input, func(key string) (string, bool) {
		switch key {
		case "VAR1":
			return "/var1", true
		case "VAR2":
			return "var2", true
		default:
			return os.LookupEnv(key)
		}
	})
	assert.Equal(t, "/var1/var2/bin", res)
}
