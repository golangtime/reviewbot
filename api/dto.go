package api

import "time"

type AddRepoRequest struct {
	Owner        string `json:"owner"`
	Name         string `json:"name"`
	MinApprovals int    `json:"min_approvals"`
	Provider     string `json:"provider"`
}

type AddRepoResponse struct {
	Success bool `json:"success"`
}

type RemoveRepoRequest struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type RemoveRepoResponse struct {
	Success bool `json:"success"`
}

type ListReposRequest struct {
	Owner string `json:"owner"`
}

type ListReposResponse struct {
	Repos   []*Repo `json:"repos"`
	Count   int     `json:"count"`
	Success bool    `json:"success"`
}

type Repo struct {
	Owner        string `json:"owner"`
	Name         string `json:"name"`
	Provider     string `json:"provider"`
	MinApprovals int    `json:"min_approvals"`
}

type ListNotificationRequest struct {
}

type ListNotificationResponse struct {
	Notifications []*Notification
}

type AddNotificationRuleRequest struct {
	UserID           int64  `json:"user_id"`
	ChatID           string `json:"chat_id"`
	NotificationType string `json:"notification_type"`
	ProviderID       string `json:"provider_id"`
	Priority         int    `json:"priority"`
}

type AddNotificationRuleResponse struct {
	Success bool `json:"success"`
}

type RemoveNotitifcationRuleRequest struct {
	ID int64 `json:"id"`
}

type RemoveNotitifcationRuleResponse struct {
	Success bool `json:"success"`
}

type ListNotificationRulesRequest struct {
}

type ListNotificationRulesResponse struct {
	Result []*NotificationRule
}

type NotificationRule struct {
	ID               int64  `json:"id"`
	UserID           int64  `json:"user_id"`
	NotificationType string `json:"notification_type"`
	ProviderID       string `json:"provider_id"`
	ChatID           int64  `json:"chat_id"`
	Priority         int    `json:"priority"`
}

type Notification struct {
	Recepient   string    `json:"recepient"`
	Link        string    `json:"link"`
	UserID      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	ReservedFor time.Time `json:"reserved_for"`
}
