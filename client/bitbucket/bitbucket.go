package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	bitbucket "github.com/gfleury/go-bitbucket-v1"

	"github.com/golangtime/reviewbot/client"
)

type BitbucketClient struct {
	g *bitbucket.APIClient
}

func New(baseUrl, username, password string) *BitbucketClient {
	basePath := baseUrl + "/rest"

	ctxAuth := context.WithValue(
		context.Background(),
		bitbucket.ContextBasicAuth,
		bitbucket.BasicAuth{
			UserName: username,
			Password: password,
		},
	)

	client := bitbucket.NewAPIClient(ctxAuth, bitbucket.NewConfiguration(basePath))

	cl := &BitbucketClient{
		g: client,
	}

	return cl
}

func (c *BitbucketClient) UnfinishedPullRequests(owner, repo string, minApprovals int) ([]*client.PullRequest, error) {
	pullRequests, err := c.ListPullRequests(owner, repo)
	if err != nil {
		return nil, err
	}

	var result []*client.PullRequest

	for _, pr := range pullRequests {
		if pr.Approvals < minApprovals {
			result = append(result, pr)
		}
	}

	return result, nil
}

func (c *BitbucketClient) ListPullRequests(owner, repoName string) ([]*client.PullRequest, error) {
	resp, err := c.g.DefaultApi.GetPullRequests(map[string]any{
		"role":  "AUTHOR",
		"state": "OPEN",
	})
	if err != nil {
		return nil, err
	}

	var result []*client.PullRequest

	for _, prData := range resp.Values["values"].([]any) {
		pr := prData.(map[string]any)
		self := pr["fromRef"].(map[string]any)

		repo := self["repository"].(map[string]any)

		projectSlug := repo["slug"].(string)
		if projectSlug != repoName {
			continue
		}

		prBody, _ := json.MarshalIndent(prData, "", "    ")

		log.Println(string(prBody))

		pullRequestID := int64(pr["id"].(float64))
		links := pr["links"].(map[string]any)
		selfLinks := links["self"].([]any)
		selfLinkData := selfLinks[0].(map[string]any)
		selfLink := selfLinkData["href"].(string)

		var approvals int
		var reviewers []client.Reviewer
		for _, rev := range pr["reviewers"].([]any) {
			r := rev.(map[string]any)

			status := r["status"].(string)

			if status != "APPROVED" {
				logData, _ := json.MarshalIndent(prData, "", "    ")
				fmt.Println(string(logData))
				userData := r["user"].(map[string]any)
				userID := int64(userData["id"].(float64))
				email := userData["emailAddress"].(string)

				reviewers = append(reviewers, client.Reviewer{
					ID:    userID,
					Email: email,
				})
			} else {
				approvals++
			}
		}

		result = append(result, &client.PullRequest{
			ExternalID: pullRequestID,
			Link:       selfLink,
			Reviewers:  reviewers,
			Approvals:  approvals,
		})
	}

	body, _ := json.MarshalIndent(result, "", "    ")
	fmt.Println(string(body))

	return result, nil
}
