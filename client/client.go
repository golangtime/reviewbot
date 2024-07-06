package client

import "github.com/google/go-github/v62/github"

type Repository struct {
}

type GitClient interface {
	ListRepositories(owner string) ([]*github.Repository, error)
	ListPullRequests(owner, repo string) ([]*github.PullRequest, error)
	// GetOpenMergeRequests(repoNames ...string)
	// // ListPeers получает список участников CodeReview для заданного MergeRequest
	// ListPeers(mergeRequest string) ([]string, error)
	// // Подсчитывает минимальное количество approval
	// CountApprovals(mergeRequest string) (int, error)
}
