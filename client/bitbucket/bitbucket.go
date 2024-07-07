package bitbucket

import (
	"fmt"
	"log"
	"net/url"

	"github.com/golangtime/reviewbot/client"
	"github.com/ktrysmt/go-bitbucket"
)

type BitbucketClient struct {
	g *bitbucket.Client
}

func New(baseUrl, username, password string) *BitbucketClient {
	client := bitbucket.NewBasicAuth(username, password)

	u, err := url.Parse(baseUrl)
	if err != nil {
		panic(err)
	}

	client.SetApiBaseURL(*u)

	res, err := client.Repositories.ListForAccount(&bitbucket.RepositoriesOptions{
		Owner: "aleksandr.nemtarev",
	})

	if err != nil {
		panic(err)
	}

	log.Println("resp", res)

	return &BitbucketClient{
		g: client,
	}
}

func (c *BitbucketClient) ListReviews(owner, repo string, id int) ([]*client.PullRequestReview, error) {
	opt := &bitbucket.PullRequestsOptions{
		ID:       fmt.Sprintf("%d", id),
		Owner:    owner,
		RepoSlug: repo,
	}

	resp, err := c.g.Repositories.PullRequests.Statuses(opt)
	if err != nil {
		return nil, err
	}

	statuses := resp.([]any)

	var result []*client.PullRequestReview

	for _, status := range statuses {
		fmt.Sprintf("type: %T, value: %+v", status, status)
		result = append(result, &client.PullRequestReview{
			Status: "-",
		})
	}

	return result, nil
}

func (c *BitbucketClient) UnfinishedPullRequests(owner, repo string, minApprovals int) ([]*client.PullRequest, error) {
	pullRequests, err := c.ListPullRequests(owner, repo)
	if err != nil {
		return nil, err
	}

	var result []*client.PullRequest

	for _, pr := range pullRequests {
		for _, u := range pr.Reviewers {
			fmt.Println("pending review", u.ID, u.Email)
		}

		reviews, err := c.ListReviews(owner, repo, int(pr.ExternalID))
		if err != nil {
			return nil, err
		}

		countPending := len(pr.Reviewers)
		for _, r := range reviews {
			if r.Status == "APPROVED" {
				countPending--
			}
		}

		if countPending < minApprovals {
			result = append(result, pr)
		}
	}

	return result, nil
}

func (c *BitbucketClient) ListRepositories(owner string) ([]*client.Repository, error) {
	opt := &bitbucket.RepositoriesOptions{Owner: "aleksandr.nemtarev@lamoda.ru"}
	result, err := c.g.Repositories.ListForAccount(opt)
	if err != nil {
		return nil, err
	}

	// for _, r := range repos {
	// 	log.Printf("Repository(id=%v,name=%v)\n", *r.ID, *r.Name)
	// }

	// TODO apply filter with DB filters
	var repos []*client.Repository

	for _, r := range result.Items {
		repos = append(repos, &client.Repository{
			Name: r.Name,
		})
	}

	return repos, nil
}

func (c *BitbucketClient) ListPullRequests(owner, repoName string) ([]*client.PullRequest, error) {
	opt := &bitbucket.PullRequestsOptions{}

	resp, err := c.g.Repositories.PullRequests.Gets(opt)
	if err != nil {
		return nil, err
	}

	pullRequests := resp.([]map[string]any)

	if len(pullRequests) == 0 {
		log.Println("repository has no pull requests")
	}

	var result []*client.PullRequest

	for _, pr := range pullRequests {
		log.Println("pull request", pr)

		var reviewers []client.Reviewer

		result = append(result, &client.PullRequest{
			Reviewers: reviewers,
		})
	}

	return result, nil
}
