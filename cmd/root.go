package cmd

import (
	"errors"
	"fmt"

	"github.com/d-led/pathdebug/common"
	"github.com/d-led/pathdebug/view"
	"github.com/spf13/cobra"
)

var (
	outputMode common.Output = common.OutputInteractive

	rootCmd = &cobra.Command{
		Use:   "pathdebug {VAR_NAME}",
		Short: "Debug path lists (non-)interactively",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return errors.New("please provide the name of the environment variable to debug")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			switch outputMode {
			case common.OutputInteractive:
				if err := view.Run(); err != nil {
					common.FailWith(fmt.Sprintf("There's been an error: %v", err.Error()))
				}
			case common.OutputTable:
				fmt.Println(view.RenderTable())
			case common.OutputJSON:
				text, err := view.RenderJson()
				if err != nil {
					common.FailWith("Rendering JSON failed: " + err.Error())
				}
				fmt.Println(string(text))
			case common.OutputCSV:
				fmt.Println(view.RenderCSV())
			default:
				common.FailWith(fmt.Sprintf("output as %v not yet implemented", outputMode))
			}
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	flags := rootCmd.Flags()
	flags.VarP(&outputMode, "output", "o", fmt.Sprintf("one of: %v", common.AllOutputs))
}

func initConfig() {
}
