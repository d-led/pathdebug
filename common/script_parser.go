package common

import (
	"strings"

	"mvdan.cc/sh/v3/syntax"
)

func ForEachVariableAssignment(key, input string, fn func(string)) {
	// parse
	r := strings.NewReader(input)
	f, err := syntax.NewParser().Parse(r, "")
	if err != nil {
		return
	}

	syntax.Walk(f, func(node syntax.Node) bool {
		switch x := node.(type) {
		case *syntax.Assign:
			if key == x.Name.Value {
				value := input[x.Value.Pos().Offset():x.Value.End().Offset()]
				fn(value)
			}
		}
		return true
	})
}
