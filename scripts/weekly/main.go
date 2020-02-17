package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// TODO: create github client, port code from https://github.com/dyweb/dy-bot/blob/master/pkg/weekly/worker.go

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
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Reconcile issue",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO: fetch issues, close old one and create new one")
		},
	}
	rootCmd.AddCommand(issueCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
