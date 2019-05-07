package main

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

func (ch *StatusChecker) ShouldAdvance(pr *PullRequest) bool {
	// Go in reverse order of comments to find if we should try to advance
	for i := len(pr.Comments.Nodes) - 1; i >= 0; i-- {
		comment := pr.Comments.Nodes[i]
		author := string(comment.Author.Login)

		switch LoadComment(string(comment.BodyText)).(type) {
		case CommentMerge:
			if ch.mergers[author] {
				return true
			}

		case CommentCancel:
			if ch.mergers[author] {
				return false
			}

		case CommentBlocked:
			return false
		}
	}

	return false
}
