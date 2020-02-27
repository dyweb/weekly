package main

import (
	"context"
	"fmt"
	"os"

	dlog "github.com/dyweb/gommon/log"
	"github.com/spf13/cobra"
)

var logReg = dlog.NewRegistry()
var log = logReg.Logger()

func main() {
	fmt.Printf("weekly args %v\n", os.Args)
	rootCmd := &cobra.Command{
		Use:   "weekly",
		Short: "dyweb weekly",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}

	var weeklyHome string
	generateCmd := &cobra.Command{
		Use:     "gen",
		Aliases: []string{"generate"},
		Short:   "Generate weekly file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return BuildRecent(context.Background(), weeklyHome)
		},
	}
	generateCmd.Flags().StringVar(&weeklyHome, "home", "", "Home of weekly directory, e.g. /home/gaocegege/workspace/src/github.com/dyweb/weekly")

	rootCmd.AddCommand(generateCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
