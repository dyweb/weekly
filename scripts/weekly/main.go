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
	var weeklyHome string
	var dryRun bool

	rootCmd := &cobra.Command{
		Use:   "weekly",
		Short: "dyweb weekly",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}
	rootCmd.PersistentFlags().StringVar(&weeklyHome, "home", "", "Home of weekly directory, e.g. /home/gaocegege/workspace/src/github.com/dyweb/weekly")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry", false, "dry run")

	generateCmd := &cobra.Command{
		Use:     "gen",
		Aliases: []string{"generate"},
		Short:   "Generate weekly file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return BuildRecent(context.Background(), weeklyHome)
		},
	}

	issueCmd := &cobra.Command{
		Use: "issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := CloseOld(context.Background(), dryRun); err != nil {
				return err
			}
			if err := OpenNew(context.Background(), dryRun); err != nil {
				return err
			}
			return nil
		},
	}

	// TODO: add issue command and by default run command in this order
	// - issue
	// - gen
	// - commit

	rootCmd.AddCommand(issueCmd)
	rootCmd.AddCommand(generateCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
