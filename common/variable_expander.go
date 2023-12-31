package common

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"go.spiff.io/expand"
)

func CustomExpandVariables(input string, fn func(key string) (value string, ok bool)) string {
	input = FixHomeExpansion(input)
	input = ConvertSimpleVarsToBraces(input)
	input = strings.TrimSpace(input)
	input = strings.Trim(input, fmt.Sprintf("%v\"'", os.PathListSeparator))
	return expand.Expand(input, fn)
}

func ConvertSimpleVarsToBraces(input string) string {
	r := regexp.MustCompile(`\$([0-9a-zA-Z_]+)`)
	return r.ReplaceAllStringFunc(input, func(m string) string {
		parts := r.FindStringSubmatch(m)
		return fmt.Sprintf(`${%s}`, parts[1])
	})
}

func FixHomeExpansion(path string) string {
	path = strings.ReplaceAll(path, "$HOME/", "~/")
	path = strings.ReplaceAll(path, "${HOME}/", "~/")
	return path
}
