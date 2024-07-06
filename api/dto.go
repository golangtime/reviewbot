package api

import "time"

type AddRepoRequest struct {
	Owner        string `json:"owner"`
	Name         string `json:"name"`
	MinApprovals int    `json:"min_approvals"`
}

type AddRepoResponse struct {
	Success bool `json:"success"`
}

type ListReposRequest struct {
	Owner string `json:"owner"`
}

type ListReposResponse struct {
	Repos []*Repo
}

type Repo struct {
	Name         string `json:"name"`
	MinApprovals int    `json:"min_approvals"`
}

type ListNotificationRequest struct {
	QueueType string `json:"queue_type"`
}

type ListNotificationResponse struct {
	Notifications []*Notification
}

type AddNotificationRuleRequest struct {
	UserID           int64  `json:"user_id"`
	NotificationType string `json:"notification_type"`
	ProviderID       string `json:"provider_id"`
	Priority         int    `json:"priority"`
}

type AddNotificationRuleResponse struct {
	Success bool `json:"success"`
}

type ListNotificationRulesRequest struct {
}

type ListNotificationRulesResponse struct {
	Result []*NotificationRule
}

type NotificationRule struct {
	UserID           int64  `json:"user_id"`
	NotificationType string `json:"notification_type"`
	ProviderID       string `json:"provider_id"`
	Priority         int    `json:"priority"`
}

type Notification struct {
	Recepient   string    `json:"recepient"`
	Link        string    `json:"link"`
	UserID      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	ReservedFor time.Time `json:"reserved_for"`
}
