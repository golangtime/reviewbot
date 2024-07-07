package github

import (
	"context"
	"log"

	"github.com/golangtime/reviewbot/client"
	"github.com/google/go-github/v62/github"
)

type GithubClient struct {
	g *github.Client
}

func New() *GithubClient {
	client := github.NewClient(nil)

	return &GithubClient{
		g: client,
	}
}

func (c *GithubClient) ListPullRequests(owner, repoName string) ([]*client.PullRequest, error) {
	pullRequests, _, err := c.g.PullRequests.List(context.Background(), owner, repoName, &github.PullRequestListOptions{})
	if err != nil {
		return nil, err
	}

	if len(pullRequests) == 0 {
		log.Println("repository has no pull requests")
	}

	var result []*client.PullRequest

	for _, pr := range pullRequests {
		var reviewers []client.Reviewer
		for _, r := range pr.RequestedReviewers {
			var email string
			if r.Email != nil {
				email = *r.Email
			}

			reviewers = append(reviewers, client.Reviewer{
				ID:    *r.ID,
				Email: email,
			})
		}

		result = append(result, &client.PullRequest{
			ExternalID: pr.GetID(),
			Link:       *pr.HTMLURL,
			Reviewers:  reviewers,
		})
	}

	return result, nil
}

func (c *GithubClient) ListReviews(owner, repoName string, prNumber int) ([]*client.PullRequestReview, error) {
	reviews, _, err := c.g.PullRequests.ListReviews(context.Background(), owner, repoName, prNumber, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []*client.PullRequestReview
	for _, r := range reviews {
		var state string
		if r.State != nil {
			state = *r.State
		}
		result = append(result, &client.PullRequestReview{
			Status: state,
		})
	}

	return result, nil
}
