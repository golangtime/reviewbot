package db

import (
	"database/sql"
	"fmt"
	"time"
)

type RepoEntity struct {
	Owner        string
	Name         string
	MinApprovals int
	// client type (github, bitbucket, etc)
	Provider string
}

type Notification struct {
	ID          int64
	Recepient   string
	Link        string
	UserID      int64
	CreatedAt   time.Time
	ReservedFor time.Time
	Status      string
}

type NotificationRule struct {
	UserID           int64
	NotificationType string
	// ProviderID chatID in messengers
	ProviderID string
	Priority   int
}

type Repo interface {
	ListRepos(db *sql.DB, owner string) ([]RepoEntity, error)
	AddRepo(db *sql.DB, owner, repo string, minApproval int, clientType string) error
	ListPendingNotifications(db *sql.DB, queueType string) ([]Notification, error)
	AddNotificationRule(db *sql.DB, userID int64, notificationType string, providerID string, priority int) error
	UpdateNotificationRule(db *sql.DB, userID int64, notificationType, providerID string, priority int) error
	ListNotificationRules(db *sql.DB) ([]*NotificationRule, error)
}

type Repository struct {
}

func (r Repository) ListRepos(db *sql.DB, owner string) ([]RepoEntity, error) {
	query := "SELECT name, owner, provider, min_approvals FROM repositories"
	var args []any
	if owner != "" {
		query += " WHERE owner = $1"
		args = append(args, owner)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var repos []RepoEntity

	for rows.Next() {
		var (
			name         string
			owner        string
			provider     string
			minApprovals int
		)

		err := rows.Scan(&name, &owner, &provider, &minApprovals)
		if err != nil {
			return nil, err
		}

		repoProvider := "github"
		if provider != "" {
			repoProvider = provider
		}

		repos = append(repos, RepoEntity{
			Name:         name,
			Owner:        owner,
			Provider:     repoProvider,
			MinApprovals: minApprovals,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}

func (r Repository) AddRepo(db *sql.DB, owner, repo string, minApproval int, provider string) error {
	_, err := db.Exec("INSERT INTO repositories (owner, name, min_approvals, provider) VALUES ($1, $2, $3, $4)",
		owner, repo, minApproval, provider)
	return err
}

func (r Repository) EnqueueNotification(db *sql.DB, sourceType string, pullRequestURL string, email string, userID int64) error {
	notificationRule, err := r.ListNotificationRuleByUser(db, userID)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err == sql.ErrNoRows {
		return nil
	}

	var tableName string
	switch notificationRule.NotificationType {
	case "email":
		tableName = "email"
	default:
		return fmt.Errorf("unknown notification type %s", notificationRule.NotificationType)
	}

	query := fmt.Sprintf(`
	INSERT INTO notification_%s_queue (
		recepient, link, user_id, created_at, reserved_for, source) 
		VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (recepient, link, user_id) DO NOTHING`, tableName)

	_, err = db.Exec(query,
		email, pullRequestURL, userID, time.Now(), time.Now().Add(time.Hour), sourceType)
	return err
}

var notificationType = map[string]struct{}{
	"email": {},
}

func (r Repository) ListPendingNotifications(db *sql.DB, queueType string) ([]Notification, error) {
	if _, ok := notificationType[queueType]; !ok {
		return nil, fmt.Errorf("invalid notification type: %v", queueType)
	}

	query := fmt.Sprintf(`SELECT id, recepient, link, user_id, created_at, reserved_for 
	FROM notification_%s_queue WHERE status = ''
	ORDER BY created_at DESC`, queueType)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []Notification

	for rows.Next() {
		var (
			id          int64
			recepient   string
			link        string
			userID      int64
			createdAt   time.Time
			reservedFor time.Time
		)

		err := rows.Scan(&id, &recepient, &link, &userID, &createdAt, &reservedFor)
		if err != nil {
			return nil, err
		}

		result = append(result, Notification{
			ID:          id,
			Recepient:   recepient,
			Link:        link,
			UserID:      userID,
			CreatedAt:   createdAt,
			ReservedFor: reservedFor,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) AddNotificationRule(db *sql.DB, userID int64, notificationType string, providerID string, priority int) error {
	switch notificationType {
	case "email":
	default:
		return fmt.Errorf("invalid notification type %s", notificationType)
	}

	query := "INSERT INTO notification_rules (user_id, notification_type, provider_id, priority) VALUES ($1, $2, $3, $4)"

	_, err := db.Exec(query, userID, notificationType, providerID, priority)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) UpdateNotificationRule(db *sql.DB, userID int64, notificationType, providerID string, priority int) error {
	switch notificationType {
	case "email":
	default:
		return fmt.Errorf("invalid notification type %s", notificationType)
	}

	_, err := db.Exec(`UPDATE notification_rules SET priority = $3, provider_id = $4
	WHERE user_id = $1 AND notification_type = $2`, userID, notificationType, priority, providerID)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) ListNotificationRuleByUser(db *sql.DB, userID int64) (*NotificationRule, error) {
	query := `SELECT provider_id, notification_type FROM notification_rules WHERE user_id = $1 ORDER BY priority LIMIT 1`

	var rule NotificationRule

	err := db.QueryRow(query, userID).Scan(&rule.ProviderID, &rule.NotificationType)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r Repository) ListNotificationRules(db *sql.DB) ([]*NotificationRule, error) {
	query := fmt.Sprintf(`SELECT user_id, notification_type, provider_id, priority FROM notification_rules`)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []*NotificationRule

	for rows.Next() {
		var (
			userID           int64
			notificationType string
			providerID       string
			priority         int
		)

		err := rows.Scan(&userID, &notificationType, &providerID, &priority)
		if err != nil {
			return nil, err
		}

		result = append(result, &NotificationRule{
			UserID:           userID,
			NotificationType: notificationType,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
