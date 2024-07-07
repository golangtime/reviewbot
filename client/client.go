package client

type Repository struct {
	Name string
}

type PullRequest struct {
	ExternalID int64
	Link       string
	Reviewers  []Reviewer
	Approvals  int
}

type Reviewer struct {
	ID    int64
	Email string
}

type PullRequestReview struct {
	UserID int64
	Status string
}

type GitClient interface {
	ListPullRequests(owner, repo string) ([]*PullRequest, error)
	UnfinishedPullRequests(owner, repo string, minApprovals int) ([]*PullRequest, error)
}
