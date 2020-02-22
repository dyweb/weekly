package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dyweb/gommon/errors"
	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/util/fsutil"
	"github.com/google/go-github/v29/github"
	"github.com/spf13/cobra"
)

var logReg = dlog.NewRegistry()
var log = logReg.Logger()

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
			//return build()
			_, err := fileNameFromIssue(context.Background(), 183)
			return err
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

// Based on https://github.com/dyweb/dy-bot/blob/master/pkg/util/weeklyutil/weekly.go
// weekly title is Weekly-157
// so @gaocegege just write strconv.Atoi(title[7:len(title)])
// It does not work properly because we don't create weekly every week, which is the case for #183
func fileNameFromIssue(ctx context.Context, id int) (string, error) {
	// TODO: move the github token logic out, wee need to create new issue as well
	client := github.NewClient(nil)
	issue, _, err := client.Issues.Get(ctx, "dyweb", "weekly", id)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(issue.GetTitle(), "Weekly-") {
		return "", errors.Errorf("Invalid weekly issue title, must start with Weekly-, got %s", issue.GetTitle())
	}
	//created := issue.GetCreatedAt()

	log.Infof("issue is %v", issue)
	// FIXME(#186): need to modify weekly generator
	return "", errors.New("not implemented, see https://github.com/dyweb/weekly/issues/186")
}
