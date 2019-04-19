package main

import "github.com/shurcooL/githubv4"

type Vars = map[string]interface{}

type Query struct {
	RateLimit  RateLimit
	Repository Repository `graphql:"repository(owner: $owner, name: $name)"`
}

type RateLimit struct {
	Remaining githubv4.Int
	Cost      githubv4.Int
}

type Repository struct {
	ID githubv4.ID

	PullRequests struct {
		Nodes []PullRequest
	} `graphql:"pullRequests(orderBy: {field: UPDATED_AT, direction: DESC}, first: $prs, states: [OPEN])"`
}

type PullRequest struct {
	ID     githubv4.ID
	URL    githubv4.URI
	Number githubv4.Int
	Body   githubv4.String

	Mergeable            githubv4.MergeableState
	PotentialMergeCommit Commit

	Commits struct {
		Nodes []struct{ Commit Commit }
	} `graphql:"commits(last: 1)"`

	Comments struct {
		Nodes []Comment
	} `graphql:"comments(last: $comments)"`
}

type Comment struct {
	BodyText githubv4.String
	Author   Author
}

type Author struct {
	Login githubv4.String
}

type Commit struct {
	OID githubv4.GitObjectID
	URL githubv4.URI

	Status struct {
		Contexts []Context
		State    githubv4.StatusState
	}

	Parents struct {
		Nodes []struct {
			OID githubv4.GitObjectID
		}
	} `graphql:"parents(first: 2)"`
}

type Context struct {
	Context githubv4.String
	State   githubv4.StatusState
}
