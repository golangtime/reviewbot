package dto

type UpdateNotificationRuleRequest struct {
	UserID           int64  `json:"user_id"`
	NotificationType string `json:"notification_type"`
	ProviderID       string `json:"provider_id"`
	Priority         int    `json:"priority"`
}

type UpdateNotificationRuleResponse struct {
	Success bool `json:"success"`
}

type ListPullRequests struct {
	Owner    string `json:"owner"`
	Repo     string `json:"repo"`
	Provider string `json:"provider"`
}

type ListPullRequestsResponse struct {
	Result []PullRequest `json:"result"`
}

type PullRequest struct {
	ID        int64                 `json:"id"`
	Link      string                `json:"link"`
	Reviewers []PullRequestReviewer `json:"reviewers"`
}

type PullRequestReviewer struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
