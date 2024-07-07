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
