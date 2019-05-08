package main

import (
	"github.com/shurcooL/githubv4"
)

//go:generate stringer -type Status

type Status int

const (
	StatusWaiting Status = iota
	StatusMergeRequested
	StatusAttemptMerge
	StatusBlocked
	StatusCanceled
)

type StatusChecker struct {
	mergers map[string]bool
}

func NewStatusChecker(cfg *Config) *StatusChecker {
	mergers := make(map[string]bool, len(cfg.Mergers))
	for _, reviewer := range cfg.Mergers {
		mergers[reviewer] = true
	}
	return &StatusChecker{mergers: mergers}
}

func (ch *StatusChecker) Status(pr *PullRequest) Status {
	var humanState, barsState stringBox
comments:
	for i := len(pr.Comments.Nodes) - 1; i >= 0; i-- {
		comment := pr.Comments.Nodes[i]
		author := string(comment.Author.Login)

		switch LoadComment(string(comment.BodyText)).(type) {
		case CommentMerge:
			if ch.mergers[author] {
				humanState.Set("merge")
			}

		case CommentCancel:
			if ch.mergers[author] {
				humanState.Set("cancel")
			}

		case CommentStarted:
			barsState.Set("started")

		case CommentBlocked:
			barsState.Set("blocked")
			break comments // ignore any previous state
		}
	}

	if state, ok := humanState.Get(); ok && state == "merge" {
		if state, ok := barsState.Get(); ok && state == "started" {
			if pr.Mergeable == githubv4.MergeableStateConflicting {
				return StatusBlocked
			}
			return StatusAttemptMerge
		}
		return StatusMergeRequested
	}
	return StatusWaiting
}

type stringBox struct {
	value string
	set   bool
}

func (o *stringBox) Clear() {
	o.value = ""
	o.set = false
}

func (o *stringBox) Set(value string) {
	if o.set {
		return
	}
	o.set = true
	o.value = value
}

func (o *stringBox) Get() (string, bool) {
	return o.value, o.set
}
