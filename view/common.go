package view

import (
	"os"

	"github.com/d-led/pathdebug/common"
)

func getResults() []common.ResultRow {
	// args validated in the root command
	source := common.NewEnvSource(os.Args[1])

	fs := &common.OsFilesystem{}

	resultsCalculator := common.NewResultsCalculator(fs, source)

	results, err := resultsCalculator.CalculateResults()
	if err != nil {
		common.FailWith(err.Error())
	}
	return results
}
