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

	var query struct {
		RateLimit struct {
			Remaining githubv4.Int
			Cost      githubv4.Int
		}

		Repository struct {
			PullRequests struct {
				Nodes []PullRequest
			} `graphql:"pullRequests(orderBy: {field: UPDATED_AT, direction: DESC}, first: $prs, states: [OPEN])"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	err = cli.Query(ctx, &query, Vars{
		"owner":    githubv4.String(cfg.Owner),
		"name":     githubv4.String(cfg.Repo),
		"prs":      githubv4.Int(1),
		"comments": githubv4.Int(0),
	})
	if err != nil {
		return errs.Wrap(err)
	}

	fmt.Printf("%+v\n", query)
	return nil
}
