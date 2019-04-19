package main

import (
	"fmt"
	"strings"
)

const (
	barsPrefix       = "Bars 🍺\n"
	headerBarsPrefix = "## " + barsPrefix
)

func LoadComment(text string) interface{} {
	if text == "r+" {
		return CommentApproval{}
	} else if text == "r-" {
		return CommentDisapproval{}
	} else if !strings.HasPrefix(text, barsPrefix) {
		return nil
	}

	text = text[len(barsPrefix):]
	fields := strings.Fields(text)

	switch {
	case strings.HasPrefix(text, "Merge "):
		return CommentMerge{
			PRNumber:  fields[len(fields)-4][1:],
			Candidate: fields[len(fields)-1],
		}

	case strings.HasPrefix(text, "Pending "):
		return CommentPending{Candidate: fields[len(fields)-1]}

	case strings.HasPrefix(text, "Starting "):
		return CommentStarting{PR: fields[len(fields)-1]}

	case strings.HasPrefix(text, "Successful "):
		return CommentStarting{PR: fields[len(fields)-1]}

	case strings.HasPrefix(text, "Failed "):
		return CommentStarting{PR: fields[len(fields)-1]}

	case strings.HasPrefix(text, "Stale "):
		return CommentStale{}

	}
	return nil
}

type CommentApproval struct{}

func (c CommentApproval) String() string { return "r+" }

type CommentDisapproval struct{}

func (c CommentDisapproval) String() string { return "r-" }

type CommentMerge struct {
	PRNumber  string
	Candidate string
}

func (c CommentMerge) String() string {
	return fmt.Sprintf(headerBarsPrefix+"Merge attempt of #%s at ref %s", c.PRNumber, c.Candidate)
}

type CommentPending struct{ Candidate string }

func (c CommentPending) String() string {
	return fmt.Sprintf(headerBarsPrefix+"Pending tests for candidate %s", c.Candidate)
}

type CommentStarting struct{ PR string }

func (c CommentStarting) String() string {
	return fmt.Sprintf(headerBarsPrefix+"Starting tests at %s", c.PR)
}

type CommentSuccess struct{ PR string }

func (c CommentSuccess) String() string {
	return fmt.Sprintf(headerBarsPrefix+"Successful tests at %s", c.PR)
}

type CommentFailure struct{ PR string }

func (c CommentFailure) String() string {
	return fmt.Sprintf(headerBarsPrefix+"Failed tests at %s", c.PR)
}

type CommentStale struct{}

func (c CommentStale) String() string {
	return headerBarsPrefix + "Stale branch: unable to test"
}
