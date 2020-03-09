package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dyweb/gommon/errors"
	"github.com/google/go-github/v29/github"
)

// issue.go reconcile weekly issues

func RecentWeeklyIssues(ctx context.Context) ([]*github.Issue, error) {
	gh, err := NewGitHub(ctx)
	if err != nil {
		return nil, err
	}
	issues, err := gh.RecentWeeklyIssues(ctx)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func CloseOld(ctx context.Context, dryRun bool) error {
	gh, err := NewGitHub(ctx)
	if err != nil {
		return err
	}

	issues, err := gh.RecentWeeklyIssues(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	var toClose []*github.Issue
	for _, issue := range issues {
		if issue.GetState() != "open" {
			continue
		}
		// Created longer than a week, must close it
		if now.Sub(issue.GetCreatedAt()) > 24*7*time.Hour {
			toClose = append(toClose, issue)
		}
	}

	if dryRun {
		for _, issue := range toClose {
			log.Infof("Should close stale issue %d %s", issue.GetNumber(), issue.GetTitle())
		}
		return nil
	}

	merr := errors.NewMultiErr()
	for _, issue := range toClose {
		log.Infof("Closing issue %d %s", issue.GetNumber(), issue.GetTitle())
		if err := gh.CloseIssue(ctx, issue.GetNumber()); err != nil {
			log.Warnf("Closing issue %d failed: %s", issue.GetNumber(), err)
			merr.Append(err)
		}
	}

	return merr.ErrorOrNil()
}

// OpenNew should be called after CloseOld
func OpenNew(ctx context.Context, dryRun bool) error {
	// TODO(at15): duplicated code
	gh, err := NewGitHub(ctx)
	if err != nil {
		return err
	}

	issues, err := gh.RecentWeeklyIssues(ctx)
	if err != nil {
		return err
	}

	maxSeq := 0
	for _, issue := range issues {
		// NOTE: you should call CloseOld before calling CreateNew.
		if issue.GetState() == "open" {
			return nil
		}
		seq, err := GetWeeklyNumber(issue)
		if err != nil {
			return err
		}
		maxSeq = max(maxSeq, seq)
	}
	if maxSeq == 0 {
		return errors.Errorf("no recent weekly issues %d", len(issues))
	}

	newSeq := maxSeq + 1
	req := FormatWeeklyIssue(newSeq)
	if dryRun {
		log.Infof("Should create issue %s", req.GetTitle())
		return nil
	}

	issue, err := gh.OpenIssue(ctx, req)
	if err != nil {
		log.Warnf("Create new issue failed: %s", err)
		return err
	}
	log.Infof("Created new issue %d %s", issue.GetNumber(), issue.GetTitle())
	return nil
}

// Copied from https://github.com/dyweb/dy-bot/blob/2fedb230d6ba21f0ebb9f1091f27d8482c331772/pkg/util/weeklyutil/weekly.go#L11
func GetWeeklyNumber(issue *github.Issue) (int, error) {
	if !IsWeeklyIssue(issue) {
		return -1, errors.Errorf("%d is not weekly issue, title: %s", issue.GetNumber(), issue.GetTitle())
	}
	title := issue.GetTitle()
	num, err := strconv.Atoi(title[7:])
	return num, err
}

func FormatWeeklyIssue(seq int) *github.IssueRequest {
	title := fmt.Sprintf("Weekly-%d", seq)
	// TODO(at15): until I have my bot working
	assignee := "at15"
	body := fmt.Sprintf("联合周报第 %d 期开始投稿 :tada:", seq)
	return &github.IssueRequest{
		Title: &title,
		Labels: &[]string{
			"working",
		},
		Assignee: &assignee,
		Body:     &body,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
