package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return build()
		},
	}
	rootCmd.AddCommand(issueCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Based on https://github.com/dyweb/dy-bot/blob/master/pkg/weekly/weekly.go#L14 buildWeekly
// shell out to weekly-gen binary from https://github.com/dyweb/dy-weekly-generator
// weekly-gen --repo dyweb/weekly --issue 183
func build() error {
	// TODO:
	issue := 183
	cmd := exec.Command("weekly-gen", "--repo", "dyweb/weekly", "--issue", strconv.Itoa(issue))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "error running weekly-gen %s", string(out))
	}
	// TODO: get weekly number based on issue
	dst := "TODO.md"
	if err = fsutil.WriteFile(dst, out); err != nil {
		return errors.Wrapf(err, "error writing weekly-gen output to %s", dst)
	}
	return nil
}
