package main

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
	"github.com/google/go-github/v29/github"
	"github.com/jinzhu/now"
)

// build.go generates weekly file name based on create time and recent issues and shells out to weekly-gen binary.
// See https://github.com/dyweb/weekly/issues/186 for background.

// BuildRecent issue get recent list of issues and generates weekly file for all of them.
func BuildRecent(ctx context.Context, home string) error {
	if home == "" {
		return errors.New("weekly home is empty")
	}

	gh, err := NewGitHub()
	if err != nil {
		return err
	}
	issues, err := gh.RecentWeeklyIssues(ctx)
	if err != nil {
		return err
	}
	// Group issues by filename, this could happen if we are testing github workflow cron or manually open and closing issues.
	issueByFile := make(map[string][]*github.Issue)
	for _, issue := range issues {
		fname := FileName(issue.GetCreatedAt())
		issueByFile[fname] = append(issueByFile[fname], issue)
	}
	for fname, issues := range issueByFile {
		p := filepath.Join(home, fname)
		if len(issues) == 1 {
			if err := BuildOne(issues[0].GetNumber(), p); err != nil {
				return err
			}
		}

		// TODO: implement merging multiple issue, open a new issue in weekly generator
	}
	return nil
}

// Based on https://github.com/dyweb/dy-bot/blob/master/pkg/weekly/weekly.go#L14 buildWeekly
// shell out to weekly-gen binary from https://github.com/dyweb/dy-weekly-generator
// weekly-gen --repo dyweb/weekly --issue 183
func BuildOne(issue int, dst string) error {
	cmd := exec.Command("weekly-gen", "--repo", "dyweb/weekly", "--issue", strconv.Itoa(issue))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "error running weekly-gen for %d %s", issue, string(out))
	}
	dir := filepath.Dir(dst)
	if err := fsutil.MkdirIfNotExists(dir); err != nil {
		return err
	}
	if err := fsutil.WriteFile(dst, out); err != nil {
		return errors.Wrapf(err, "error writing weekly-gen output to %s", dst)
	}
	return nil
}

// FileName generates file name based on issue creation time.
// It does not rely on weekly id like https://github.com/dyweb/dy-bot/blob/master/pkg/util/weeklyutil/weekly.go
// Because it is no longer consistent.
// All the issues that are created within a week are aligned to that week's Monday.
// The comments from multiple issues that shares a same file name should be merged.
func FileName(createTime time.Time) string {
	// Copied from https://github.com/dyweb/dy-bot/blob/2fedb230d6ba21f0ebb9f1091f27d8482c331772/pkg/util/weeklyutil/weekly.go#L11
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)

	cfg := now.Config{
		WeekStartDay: time.Monday,
		TimeLocation: beijing,
	}
	date := cfg.With(createTime).BeginningOfWeek()
	return fmt.Sprintf("%d/%d-%.2d-%.2d-weekly.md", date.Year(), date.Year(), int(date.Month()), date.Day())
}
