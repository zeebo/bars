package main

import (
	"fmt"
	"strings"
)

const barsPrefix = "Bars 🍺\n"

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
	return fmt.Sprintf("## Bars :beer:\nMerge attempt of #%s at ref %s", c.PRNumber, c.Candidate)
}

type CommentPending struct{ Candidate string }

func (c CommentPending) String() string {
	return fmt.Sprintf("## Bars :beer:\nPending tests for candidate %s", c.Candidate)
}

type CommentStarting struct{ PR string }

func (c CommentStarting) String() string {
	return fmt.Sprintf("## Bars :beer:\nStarting tests at %s", c.PR)
}

type CommentSuccess struct{ PR string }

func (c CommentSuccess) String() string {
	return fmt.Sprintf("## Bars :beer:\nSuccessful tests at %s", c.PR)
}

type CommentFailure struct{ PR string }

func (c CommentFailure) String() string {
	return fmt.Sprintf("## Bars :beer:\nFailed tests at %s", c.PR)
}

type CommentStale struct{}

func (c CommentStale) String() string {
	return "## Bars :beer:\nStale branch: unable to test"
}
