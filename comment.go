package main

import (
	"strings"
)

const (
	barsPrefix       = "Bars 🍺\n"
	headerBarsPrefix = "## " + barsPrefix
)

func LoadComment(text string) interface{} {
	if text == "merge" {
		return CommentMerge{}
	} else if text == "cancel" {
		return CommentCancel{}
	} else if !strings.HasPrefix(text, barsPrefix) {
		return nil
	}

	text = text[len(barsPrefix):]

	switch {
	case strings.HasPrefix(text, "Merge requested."):
		return CommentStarted{}

	case strings.HasPrefix(text, "Blocked: "):
		return CommentBlocked{Why: text[len("Blocked: "):]}
	}

	return nil
}

type CommentMerge struct{}

func (c CommentMerge) String() string { return "merge" }

type CommentCancel struct{}

func (c CommentCancel) String() string { return "cancel" }

type CommentStarted struct{}

func (c CommentStarted) String() string {
	return headerBarsPrefix + "Merge requested."
}

type CommentBlocked struct {
	Why string
}

func (c CommentBlocked) String() string {
	return headerBarsPrefix + "Blocked: " + c.Why
}
