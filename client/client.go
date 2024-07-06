package client

import "github.com/google/go-github/v62/github"

type Repository struct {
}

type GitClient interface {
	ListRepositories(owner string) ([]*github.Repository, error)
	ListPullRequests(owner, repo string) ([]*github.PullRequest, error)
	ListReviewers(owner, repo string, number int) ([]*github.User, error)
	ListReviews(owner, repo string, number int) ([]*github.PullRequestReview, error)
}
