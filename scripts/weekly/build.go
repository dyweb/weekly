package main

import (
	"bytes"
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

	// Call dy-weekly-generator
	for fname, issues := range issueByFile {
		p := filepath.Join(home, fname)
		if len(issues) == 1 {
			if err := BuildOne(issues[0].GetNumber(), p); err != nil {
				return err
			}
		}

		var ids []int
		for _, issue := range issues {
			ids = append(ids, issue.GetNumber())
		}
		if err := MergeIssues(ids, p); err != nil {
			return err
		}
	}
	return nil
}

// Based on https://github.com/dyweb/dy-bot/blob/master/pkg/weekly/weekly.go#L14 buildWeekly
// shell out to weekly-gen binary from https://github.com/dyweb/dy-weekly-generator
// weekly-gen --repo dyweb/weekly --issue 183
func BuildOne(issue int, dst string) error {
	out, err := RunGenerator(issue)
	if err != nil {
		return err
	}
	return WriteWeeklyFile(out, dst)
}

const defaultFrontMatter = `---
layout: post
title: Weekly
category: Weekly
author: 东岳
---`

// MergeIssues write multiple issues into one weekly file.
func MergeIssues(issueIds []int, dst string) error {
	merged := []byte(defaultFrontMatter)
	for _, id := range issueIds {
		out, err := RunGenerator(id)
		if err != nil {
			return err
		}
		out = bytes.Trim(out, "\n")
		out = StripHeader(out)
		merged = append(merged, '\n')
		merged = append(merged, out...)
		merged = append(merged, '\n')
	}
	return WriteWeeklyFile(merged, dst)
}

// TODO(at15): strip header flag https://github.com/dyweb/dy-weekly-generator/issues/25
func RunGenerator(issue int) ([]byte, error) {
	cmd := exec.Command("weekly-gen", "--repo", "dyweb/weekly", "--issue", strconv.Itoa(issue))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "error running weekly-gen for %d %s", issue, string(out))
	}
	return out, nil
}

func WriteWeeklyFile(b []byte, dst string) error {
	dir := filepath.Dir(dst)
	if err := fsutil.MkdirIfNotExists(dir); err != nil {
		return err
	}
	if err := fsutil.WriteFile(dst, b); err != nil {
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

var frontMatterDash = []byte("---")

// StripHeader remove the front matter https://jekyllrb.com/docs/front-matter/ from generated weekly.
// TODO: work around until https://github.com/dyweb/dy-weekly-generator/issues/25 is fixed
func StripHeader(b []byte) []byte {
	if !bytes.HasPrefix(b, frontMatterDash) {
		return b
	}
	b = b[len(frontMatterDash):]
	end := bytes.Index(b, frontMatterDash)
	if end == -1 {
		return b
	}
	return b[end+len(frontMatterDash):]
}
