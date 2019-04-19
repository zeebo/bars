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
	cli := githubv4.NewClient(hcli)

	// Query the information we need to advance the state.
	var query Query
	err = cli.Query(ctx, &query, Vars{
		"owner":    githubv4.String(cfg.Owner),
		"name":     githubv4.String(cfg.Repo),
		"prs":      githubv4.Int(99),
		"comments": githubv4.Int(100),
		// "masterRef": githubv4.String(cfg.MasterRef),
	})
	if err != nil {
		return errs.Wrap(err)
	}

	fmt.Printf("rate limit: used:%v remaining:%v\n", query.RateLimit.Cost, query.RateLimit.Remaining)
	prs := query.Repository.PullRequests.Nodes

	// Compute the status of every PR.
	checker := NewStatusChecker(cfg)
	statuses := make([]Status, len(prs))
	for i := range prs {
		pr := &prs[i]

		status := checker.Status(pr)
		fmt.Printf("#%v: current status: %v\n", pr.Number, status)
		statuses[i] = status
	}

	// Try to find any merge attempt PR.
	var activeMerge *CommentMerge
	for i := range prs {
		pr := &prs[i]
		activeMerge, _ = LoadComment(string(pr.Body)).(*CommentMerge)
		if activeMerge == nil {
			continue
		}

		if pullRequestBad(pr) {
			comment := CommentFailure{PR: pr.URL.String()}
			err := addComment(ctx, cli, activeMerge.PRNumber, comment.String())
			fmt.Printf("#%v: adding failure comment: %v\n", activeMerge.PRNumber, err)

			err = closePullRequest(ctx, cli, pr.ID)
			fmt.Printf("#%v: closing pull request: %v\n", pr.Number, err)

		} else if pullRequestMergable(pr) {
			comment := CommentSuccess{PR: pr.URL.String()}
			err := addComment(ctx, cli, activeMerge.PRNumber, comment.String())
			fmt.Printf("#%v: adding success comment: %v\n", activeMerge.PRNumber, err)

			err = closePullRequest(ctx, cli, pr.ID)
			fmt.Printf("#%v: closing pull request: %v\n", pr.Number, err)
		}

		break
	}

	// Update every PR to the next state.
	for i := range prs {
		pr := &prs[i]
		commits := pr.Commits.Nodes
		lastCommit := commits[len(commits)-1].Commit
		status := statuses[i]

		switch status {
		case StatusApproved:
			comment := CommentPending{Candidate: string(lastCommit.OID)}
			err := addComment(ctx, cli, pr.ID, comment.String())
			fmt.Printf("#%v: adding pending comment: %v\n", pr.Number, err)

		case StatusPending:
			if activeMerge != nil {
				break
			}

			activeMerge = &CommentMerge{
				PRNumber:  fmt.Sprint(pr.Number),
				Candidate: string(lastCommit.OID),
			}

			testPR, err := createPullRequest(ctx, cli, CreatePullRequestInput{
				BaseRefName:  githubv4.String(cfg.MasterRef),
				Body:         githubv4.String(activeMerge.String()),
				HeadRefName:  githubv4.String(pr.PotentialMergeCommit.OID),
				RepositoryId: query.Repository.ID,
				Title:        "Bars: Testing PR",
			})
			fmt.Printf("#%v: opening pull request: %v\n", pr.Number, err)

			if err == nil {
				comment := CommentStarting{PR: testPR.URL.String()}
				err := addComment(ctx, cli, pr.ID, comment.String())
				fmt.Printf("#%v: adding starting comment: %v\n", pr.Number, err)
			}

			err = closePullRequest(ctx, cli, pr.ID)
			fmt.Printf("#%v: closing pull request: %v\n", pr.Number, err)

		case StatusStale:
			comment := CommentStale{}
			err := addComment(ctx, cli, pr.ID, comment.String())
			fmt.Printf("#%v: adding stale comment: %v\n", pr.Number, err)

		case StatusTested:
			// TODO: merge pr
		}
	}
	return nil
}
