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
	})
	if err != nil {
		return errs.Wrap(err)
	}

	fmt.Printf("rate limit: used:%v remaining:%v\n", query.RateLimit.Cost, query.RateLimit.Remaining)
	prs := query.Repository.PullRequests.Nodes
	checker := NewStatusChecker(cfg)

	for i := range prs {
		pr := &prs[i]

		fmt.Printf("%d: mergable:%v status:%v advance:%v\n",
			pr.Number, pr.Mergeable, pr.MergeStateStatus,
			checker.ShouldAdvance(pr))
	}

	return nil
}
