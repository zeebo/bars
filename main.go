package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shurcooL/githubv4"
	"github.com/zeebo/errs"
	"github.com/zeebo/mon"
	"golang.org/x/oauth2"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("%+v\n", err)
	}
}

func run(ctx context.Context) (err error) {
	defer mon.Start().Stop(&err)

	// Construct our github API client.
	if len(os.Args) < 2 {
		return errs.New("usage: bars <path to config>")
	}
	cfg, err := LoadConfig(os.Args[1])
	if err != nil {
		return err
	}
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
	hcli := oauth2.NewClient(context.Background(), src)
	hcli.Transport = &AcceptRoundTripper{
		RoundTripper: hcli.Transport,
		Accept:       []string{"application/vnd.github.merge-info-preview+json"},
	}
	cli := githubv4.NewClient(hcli)

	// Query the information we need to advance the state.
	var query Query
	err = cli.Query(ctx, &query, Vars{
		"owner":    githubv4.String(cfg.Owner),
		"name":     githubv4.String(cfg.Repo),
		"prs":      githubv4.Int(99),
		"comments": githubv4.Int(100),
	})
	if err != nil {
		return errs.Wrap(err)
	}

	fmt.Printf("rate limit: used:%v remaining:%v\n", query.RateLimit.Cost, query.RateLimit.Remaining)

	prs := query.Repository.PullRequests.Nodes

	// Calculate the status of every PR
	checker := NewStatusChecker(cfg)
	statuses := make(map[int]Status)
	attemptRequested := false
	for i := range prs {
		status := checker.Status(&prs[i])
		statuses[i] = status
		attemptRequested = attemptRequested || status == StatusAttemptMerge
	}

prs:
	for i := range prs {
		pr := &prs[i]
		switch statuses[i] {
		case StatusMergeRequested:
			// If we're already attempting a merge somewhere else, do nothing
			if attemptRequested {
				fmt.Printf("%v: waiting for other pr to finish before starting", pr.URL)
				continue prs
			}

			// Comment to inform that we're going to attempt to merge this
			comment := CommentStarted{}
			err := addComment(ctx, cli, pr.ID, comment.String())
			fmt.Printf("%v: adding started comment: %v\n", pr.URL, err)

		case StatusWaiting:
			fmt.Printf("%v: waiting for human to request merge\n", pr.URL)

		case StatusBlocked:
			// Comment to inform that a human is required
			comment := CommentBlocked{Why: "unknown"}
			err := addComment(ctx, cli, pr.ID, comment.String())
			fmt.Printf("%v: adding blocked comment: %v\n", pr.URL, err)

		case StatusAttemptMerge:
			fmt.Printf("%v: must merge this somehow\n", pr.URL)
		}
	}

	return nil
}
