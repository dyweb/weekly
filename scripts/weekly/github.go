package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/dyweb/gommon/util/httputil"
	"golang.org/x/oauth2"

	"github.com/google/go-github/v29/github"
)

var cachedGh *GitHub

// GitHub is a wrapper for go-github
type GitHub struct {
	client *github.Client
}

func NewGitHub(ctx context.Context) (*GitHub, error) {
	if cachedGh != nil {
		return cachedGh, nil
	}
	hc := httputil.NewUnPooledClient()
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
		})
		// TODO(at15): should save existing http client in context so the transport can be reused
		hc = oauth2.NewClient(ctx, ts)
	}
	client := github.NewClient(hc)
	gh := &GitHub{client: client}
	cachedGh = gh
	return gh, nil
}

func IsWeeklyIssue(issue *github.Issue) bool {
	return strings.HasPrefix(issue.GetTitle(), "Weekly-")
}

func (g *GitHub) RecentWeeklyIssues(ctx context.Context) ([]*github.Issue, error) {
	issues, _, err := g.client.Issues.ListByRepo(ctx, "dyweb", "weekly", &github.IssueListByRepoOptions{
		State:     "all",
		Sort:      "created",
		Direction: "desc",
		Since:     time.Now().Add(-1 * 2 * 30 * 24 * time.Hour),
	})
	if err != nil {
		return nil, err
	}
	var filtered []*github.Issue
	for _, issue := range issues {
		if IsWeeklyIssue(issue) {
			filtered = append(filtered, issue)
		}
		//log.Infof("issue %d %d %s", issue.GetID(), issue.GetNumber(), issue.GetTitle())
	}
	return filtered, nil
}

func (g *GitHub) Issue(ctx context.Context, id int) (*github.Issue, error) {
	issue, _, err := g.client.Issues.Get(ctx, "dyweb", "weekly", id)
	return issue, err
}

func (g *GitHub) OpenIssue(ctx context.Context, req *github.IssueRequest) (*github.Issue, error) {
	issue, _, err := g.client.Issues.Create(ctx, "dyweb", "weekly", req)
	return issue, err
}

// CloseIssue close and remove ALL labels
func (g *GitHub) CloseIssue(ctx context.Context, id int) error {
	s := "closed"
	// TODO(at15): should only remove the working label
	var noLabel []string
	req := &github.IssueRequest{
		State:  &s,
		Labels: &noLabel,
	}
	_, _, err := g.client.Issues.Edit(ctx, "dyweb", "weekly", id, req)
	return err
}
