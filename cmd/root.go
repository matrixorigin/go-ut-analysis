package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func initRootCommand() *cobra.Command {
	return &cobra.Command{
		Use: "go-ut-analysis",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}

func Execute() {
	rootCommand := initRootCommand()
	// add sub command to root command
	rootCommand.AddCommand(initTestCommand())

	// run
	err := rootCommand.Execute()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
