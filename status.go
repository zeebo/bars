package main

import (
	"github.com/shurcooL/githubv4"
)

//go:generate stringer -type Status

type Status int

const (
	StatusUnknown Status = iota
	StatusUnreviewed
	StatusApproved // action: comment that it's pending
	StatusDisapproved
	StatusStale   // action: comment that it's stale
	StatusPending // action: try to make a pr if available
	StatusTesting
	StatusTested // action: land the change
	StatusMergeAttempt
)

type StatusChecker struct {
	reviewers map[string]bool
}

func NewStatusChecker(cfg *Config) *StatusChecker {
	reviewers := make(map[string]bool, len(cfg.Reviewers))
	for _, reviewer := range cfg.Reviewers {
		reviewers[reviewer] = true
	}
	return &StatusChecker{reviewers: reviewers}
}

func (ch *StatusChecker) Status(pr *PullRequest) Status {
	// Go in reverse order of comments to find current state
	reviews := make(map[string]bool)
	approvals, disapprovals := 0, 0

comments:
	for i := len(pr.Comments.Nodes) - 1; i >= 0; i-- {
		comment := pr.Comments.Nodes[i]
		author := string(comment.Author.Login)
		switch LoadComment(string(comment.BodyText)).(type) {
		case CommentApproval:
			if _, ok := reviews[author]; !ok && ch.reviewers[author] {
				reviews[author] = true
				approvals++
			}

		case CommentDisapproval:
			if _, ok := reviews[author]; !ok && ch.reviewers[author] {
				reviews[author] = false
				disapprovals++
			}

		case CommentPending:
			// We're in a pending state so ensure it's not stale
			switch pr.Mergeable {
			case githubv4.MergeableStateConflicting:
				return StatusStale
			case githubv4.MergeableStateUnknown:
				return StatusUnknown
			}
			return StatusPending

		case CommentStarting:
			return StatusTesting

		case CommentSuccess:
			// If there's a disapproval even after it's tested, then it's disapproved
			if disapprovals > 0 {
				return StatusDisapproved
			}
			return StatusTested // approved and tested

		// In these states, we essentially reset. Check current approval status.
		case CommentStale, CommentFailure:
			break comments
		}
	}

	// At this point, we're either after a failure or before any tests.
	// First, check if we're in a disapproved or unreviewed state.
	if disapprovals > 0 {
		return StatusDisapproved
	} else if approvals < 1 {
		return StatusUnreviewed
	}

	// We're in an approved state so ensure it's not stale.
	switch pr.Mergeable {
	case githubv4.MergeableStateConflicting:
		return StatusStale
	case githubv4.MergeableStateUnknown:
		return StatusUnknown
	}
	return StatusApproved
}

// pullRequestMergable returns true if the pull request should be merged.
func pullRequestMergable(pr *PullRequest) bool {
	if pr.Mergeable != githubv4.MergeableStateMergeable {
		return false
	}

	for _, commit := range pr.Commits.Nodes {
		for _, context := range commit.Commit.Status.Contexts {
			if context.State != githubv4.StatusStateSuccess {
				return false
			}
		}
	}

	return true
}

// pullRequestBad returns true if the pull request shouldn't be merged.
func pullRequestBad(pr *PullRequest) bool {
	if pr.Mergeable == githubv4.MergeableStateConflicting {
		return true
	}

	for _, commit := range pr.Commits.Nodes {
		for _, context := range commit.Commit.Status.Contexts {
			if context.State != githubv4.StatusStateSuccess &&
				context.State != githubv4.StatusStatePending {
				return true
			}
		}
	}

	return false
}
