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

type MergeStateStatus string

const (
	MergeStateStatusBehind   MergeStateStatus = "BEHIND"
	MergeStateStatusBlocked  MergeStateStatus = "BLOCKED"
	MergeStateStatusClean    MergeStateStatus = "CLEAN"
	MergeStateStatusDirty    MergeStateStatus = "DIRTY"
	MergeStateStatusDraft    MergeStateStatus = "DRAFT"
	MergeStateStatusHasHooks MergeStateStatus = "HAS_HOOKS"
	MergeStateStatusUnknown  MergeStateStatus = "UNKNOWN"
	MergeStateStatusUnstable MergeStateStatus = "UNSTABLE"
)

type PullRequest struct {
	ID               githubv4.ID
	Number           githubv4.Int
	URL              githubv4.URI
	State            githubv4.PullRequestState
	Mergeable        githubv4.MergeableState
	MergeStateStatus MergeStateStatus

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
