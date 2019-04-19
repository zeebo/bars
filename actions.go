package main

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/zeebo/errs"
)

func addComment(ctx context.Context, cli *githubv4.Client, subject githubv4.ID, body string) error {
	var m struct {
		AddComment struct {
			ClientMutationId githubv4.String
		} `graphql:"addComment(input: $input)"`
	}
	input := githubv4.AddCommentInput{
		SubjectID: subject,
		Body:      githubv4.String(body),
	}
	return errs.Wrap(cli.Mutate(ctx, &m, input, nil))
}

type CreatePullRequestInput struct {
	BaseRefName  githubv4.String `json:"baseRefName"`
	Body         githubv4.String `json:"body"`
	HeadRefName  githubv4.String `json:"headRefName"`
	RepositoryId githubv4.ID     `json:"repositoryId"`
	Title        githubv4.String `json:"title"`
}

func createPullRequest(ctx context.Context, cli *githubv4.Client, input CreatePullRequestInput) (*PullRequest, error) {
	var m struct {
		CreatePullRequest struct {
			PullRequest PullRequest
		} `graphql:"createPullRequest(input: $input)"`
	}
	err := errs.Wrap(cli.Mutate(ctx, &m, input, Vars{
		"comments": githubv4.Int(0),
	}))
	if err != nil {
		return nil, err
	}
	return &m.CreatePullRequest.PullRequest, nil
}

func closePullRequest(ctx context.Context, cli *githubv4.Client, id githubv4.ID) error {
	var m struct {
		ClosePullRequest struct {
			ClientMutationId githubv4.String
		} `graphql:"closePullRequest(input: $input)"`
	}
	type ClosePullRequestInput struct {
		PullRequestID githubv4.ID `json:"pullRequestId"`
	}
	input := ClosePullRequestInput{PullRequestID: id}
	return errs.Wrap(cli.Mutate(ctx, &m, input, nil))
}
