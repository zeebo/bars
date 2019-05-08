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
