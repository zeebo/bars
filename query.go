package main

import "github.com/shurcooL/githubv4"

type Vars = map[string]interface{}

type PullRequest struct {
	URL                  githubv4.URI
	Mergeable            githubv4.MergeableState
	PotentialMergeCommit Commit
	Commits              struct {
		Nodes []struct {
			Commit Commit
		}
	} `graphql:"commits(last: 1)"`
	Comments struct {
		Nodes []struct {
			BodyText githubv4.String
			Author   struct{ Login githubv4.String }
		}
	} `graphql:"comments(last: $comments)"`
}

type Commit struct {
	OID    githubv4.GitObjectID
	URL    githubv4.URI
	Status struct {
		Contexts []struct {
			Context githubv4.String
			State   githubv4.StatusState
		}
		State githubv4.StatusState
	}
}
