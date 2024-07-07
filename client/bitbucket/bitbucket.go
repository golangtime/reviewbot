package bitbucket

import (
	"log"

	"github.com/golangtime/reviewbot/client"
	"github.com/google/go-github/v62/github"
	"github.com/ktrysmt/go-bitbucket"
)

type BitbucketClient struct {
	g *bitbucket.Client
}

func New() *BitbucketClient {
	client := bitbucket.NewBasicAuth("username", "password")

	return &BitbucketClient{
		g: client,
	}
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

	// TODO apply filter with DB filters

	return result, nil
}

func (c *BitbucketClient) ListReviews(owner, repoName string, prNumber int) ([]*github.PullRequestReview, error) {
	return nil, nil
}

func (c *BitbucketClient) ListReviewers(owner, repoName string, prNumber int) ([]*github.User, error) {
	return nil, nil
}
