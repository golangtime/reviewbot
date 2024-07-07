package client

type Repository struct {
	Name string
}

type PullRequest struct {
	ExternalID int64
	Link       string
	Reviewers  []Reviewer
}

type Reviewer struct {
	ID    int64
	Email string
}

type PullRequestReview struct {
	Status string
}

type GitClient interface {
	ListPullRequests(owner, repo string) ([]*PullRequest, error)
	ListReviews(owner, repo string, number int) ([]*PullRequestReview, error)
}
