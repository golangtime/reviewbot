package github

import (
	"context"
	"fmt"
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

func (c *GithubClient) UnfinishedPullRequests(owner, repo string, minApprovals int) ([]*client.PullRequest, error) {
	pullRequests, err := c.ListPullRequests(owner, repo)
	if err != nil {
		return nil, err
	}

	var result []*client.PullRequest

	for _, pr := range pullRequests {
		reviewers := map[int64]struct{}{}

		for _, u := range pr.Reviewers {
			fmt.Println("pending review", u.ID, u.Email)
			reviewers[u.ID] = struct{}{}
		}

		reviews, err := c.ListReviews(owner, repo, int(pr.ExternalID))
		if err != nil {
			return nil, err
		}

		approvedCount := 0
		for _, r := range reviews {
			if r.Status == "APPROVED" {
				approvedCount++
				delete(reviewers, r.UserID)
			}
		}

		log.Println("repo approve count", approvedCount)
		log.Println("repo need notify", approvedCount < minApprovals)

		if approvedCount < minApprovals {
			result = append(result, pr)
		}
	}

	return result, nil
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

		log.Printf("pull request (%d), link: %v, reviewers: %d\n", pr.GetNumber(), pr.HTMLURL, len(reviewers))

		result = append(result, &client.PullRequest{
			ExternalID: int64(pr.GetNumber()),
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

		userID := *r.User.ID

		result = append(result, &client.PullRequestReview{
			Status: state,
			UserID: userID,
		})
	}

	return result, nil
}
