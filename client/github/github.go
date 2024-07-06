package github

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/go-github/v62/github"
)

type GithubClient struct {
	g *github.Client
}

func New() *GithubClient {
	client := github.NewClient(nil) //.
	// WithAuthToken("")

	return &GithubClient{
		g: client,
	}
}

func (c *GithubClient) ListRepositories(owner string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByUserOptions{Type: "public"}
	repos, _, err := c.g.Repositories.ListByUser(context.Background(), owner, opt)
	if err != nil {
		return nil, err
	}

	// for _, r := range repos {
	// 	log.Printf("Repository(id=%v,name=%v)\n", *r.ID, *r.Name)
	// }

	// TODO apply filter with DB filters

	return repos, nil
}

func (c *GithubClient) ListPullRequests(owner, repoName string) ([]*github.PullRequest, error) {
	pullRequests, _, err := c.g.PullRequests.List(context.Background(), owner, repoName, &github.PullRequestListOptions{})
	if err != nil {
		return nil, err
	}

	// for _, pr := range pullRequests {
	// 	body, err := json.MarshalIndent(pr, "", "    ")
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	log.Println(string(body))
	// }

	if len(pullRequests) == 0 {
		log.Println("repository has no pull requests")
	}

	// TODO apply filter with DB filters

	return pullRequests, nil
}

func (c *GithubClient) ListReviews(owner, repoName string, prNumber int) ([]*github.PullRequestReview, error) {
	reviews, _, err := c.g.PullRequests.ListReviews(context.Background(), owner, repoName, prNumber, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (c *GithubClient) ListReviewers(owner, repoName string, prNumber int) ([]*github.User, error) {
	reviewers, _, err := c.g.PullRequests.ListReviewers(context.Background(), owner, repoName, 1, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pr := range reviewers.Users {
		body, err := json.MarshalIndent(pr, "", "    ")
		if err != nil {
			return nil, err
		}
		log.Println(string(body))
	}

	if len(reviewers.Users) == 0 {
		log.Println("pull has no reviewers")
	}

	// TODO apply filter with DB filters

	return reviewers.Users, nil
}
