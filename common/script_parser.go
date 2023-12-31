package common

import (
	"strings"

	"mvdan.cc/sh/v3/syntax"
)

func ForEachVariableAssignment(key, input string, fn func(string)) {
	f, err := parseScript(input)
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

var isSourceCommand = map[string]bool{
	".":      true,
	"source": true,
	// eval?
}

func ForEachSourcedScript(input string, fn func(string)) {
	f, err := parseScript(input)
	if err != nil {
		return
	}

	syntax.Walk(f, func(node syntax.Node) bool {
		switch x := node.(type) {
		case *syntax.CallExpr:
			if len(x.Args) == 2 {
				cmd := input[x.Args[0].Pos().Offset():x.Args[0].End().Offset()]
				cmd = strings.TrimLeft(cmd, "\\")

				arg := input[x.Args[1].Pos().Offset():x.Args[1].End().Offset()]
				arg = strings.Trim(arg, "\"'`")

				if !isSourceCommand[cmd] {
					return true
				}

				fn(arg)
			}
		}
		return true
	})
}

func parseScript(input string) (*syntax.File, error) {
	r := strings.NewReader(input)
	f, err := syntax.NewParser().Parse(r, "")
	return f, err
}
