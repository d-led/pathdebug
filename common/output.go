package common

import (
	"fmt"
	"strings"
)

type Output string

const (
	OutputInteractive Output = "interactive" // default
	OutputTable       Output = "table"
	OutputJSON        Output = "json"
	OutputCSV         Output = "csv"
)

var AllOutputs []string = []string{
	string(OutputInteractive),
	string(OutputTable),
	string(OutputJSON),
	string(OutputCSV),
}

func (e *Output) String() string {
	return string(*e)
}

func (e *Output) Set(v string) error {
	switch v {
	case string(OutputInteractive), string(OutputJSON), string(OutputTable), string(OutputCSV):
		*e = Output(v)
		return nil
	default:
		return fmt.Errorf(`must be one of %v, got %v`, strings.Join(AllOutputs, ", "), v)
	}
}

func (e *Output) Type() string {
	return "output"
}
