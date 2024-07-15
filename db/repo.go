package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

var (
	ErrNotificationRuleNotFound = errors.New("notification rule not found")
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
	SenderType  string
}

type NotificationRule struct {
	ID               int64
	UserID           int64
	ChatID           int64
	NotificationType string
	// ProviderID chatID in messengers
	ProviderID string
	Priority   int
}

type PullRequestEntity struct {
	ID        int64
	Link      string
	Reviewers []ReviewEntity
}

type ReviewEntity struct {
	ID int64
}

type Repo interface {
	ListRepos(db *sql.DB, owner string) ([]RepoEntity, error)
	AddRepo(db *sql.DB, owner, repo string, minApproval int, clientType string) error
	RemoveRepo(db *sql.DB, owner, repo string) error
	ListPendingNotifications(db *sql.DB) ([]Notification, error)
	AddNotificationRule(db *sql.DB, userID int64, notificationType string, providerID, chatID string, priority int) error
	UpdateNotificationRule(db *sql.DB, userID int64, notificationType, providerID string, priority int) error
	RemoveNotificationRule(db *sql.DB, ruleID int64) error
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

func (r Repository) RemoveRepo(db *sql.DB, owner, repo string) error {
	_, err := db.Exec("DELETE FROM repositories WHERE owner = $1 AND name = $2",
		owner, repo)
	return err
}

func (r Repository) EnqueueNotification(db *sql.DB, sourceType string, pullRequestURL string, email string, userID int64) error {
	notificationRule, err := r.ListNotificationRuleByUser(db, userID)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err == sql.ErrNoRows {
		log.Println("skip enqueue notification, not notification rules", userID)
		return nil
	}

	query := fmt.Sprintf(`
	INSERT INTO notification_queue (
		rule_id, recepient, link, user_id, created_at, reserved_for, source) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (rule_id, recepient, link, user_id) DO NOTHING`)

	_, err = db.Exec(query,
		notificationRule.ID, email, pullRequestURL, userID, time.Now(), time.Now().Add(time.Hour), sourceType)
	return err
}

var notificationType = map[string]struct{}{
	"email": {},
}

func (r Repository) ListPendingNotifications(db *sql.DB) ([]Notification, error) {
	query := fmt.Sprintf(`SELECT rule_id, id, recepient, link, user_id, created_at, reserved_for 
	FROM notification_queue WHERE status = ''
	ORDER BY created_at DESC`)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []Notification

	for rows.Next() {
		var (
			ruleID      int64
			id          int64
			recepient   string
			link        string
			userID      int64
			createdAt   time.Time
			reservedFor time.Time
		)

		err := rows.Scan(&ruleID, &id, &recepient, &link, &userID, &createdAt, &reservedFor)
		if err != nil {
			return nil, err
		}

		notificationRule, err := r.GetNotificationRule(db, ruleID)
		if err != nil {
			if errors.Is(err, ErrNotificationRuleNotFound) {
				continue
			}
			return nil, err
		}

		result = append(result, Notification{
			ID:          id,
			Recepient:   recepient,
			Link:        link,
			UserID:      userID,
			CreatedAt:   createdAt,
			ReservedFor: reservedFor,
			SenderType:  notificationRule.NotificationType,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) GetNotificationRule(db *sql.DB, ruleID int64) (*NotificationRule, error) {
	var notificationRule NotificationRule
	err := db.QueryRow(`SELECT id, user_id, notification_type, provider_id, coalesce(chat_id, '0'), priority
		FROM notification_rules WHERE id = $1`, ruleID).
		Scan(&notificationRule.ID,
			&notificationRule.UserID,
			&notificationRule.NotificationType,
			&notificationRule.ProviderID,
			&notificationRule.ChatID,
			&notificationRule.Priority,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotificationRuleNotFound
		}
		return nil, fmt.Errorf("error getting notification rule with id = %d: %w", ruleID, err)
	}

	return &notificationRule, nil
}

func (r Repository) AddNotificationRule(db *sql.DB, userID int64, notificationType string, providerID, chatID string, priority int) error {
	switch notificationType {
	case "email", "pachca":
	default:
		return fmt.Errorf("invalid notification type %s", notificationType)
	}

	var chatIDVal *string
	if chatID != "" {
		chatIDVal = &chatID
	}

	query := "INSERT INTO notification_rules (user_id, notification_type, provider_id, chat_id, priority) VALUES ($1, $2, $3, $4, $5)"

	_, err := db.Exec(query, userID, notificationType, providerID, chatIDVal, priority)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) UpdateNotificationRule(db *sql.DB, userID int64, notificationType, providerID string, priority int) error {
	switch notificationType {
	case "email", "pachca":
	default:
		return fmt.Errorf("invalid notification type %s", notificationType)
	}

	log.Printf("update notification rule, user_id=%d, notification_type=%s, provider_id=%s, priority=%d\n",
		userID, notificationType, providerID, priority)

	res, err := db.Exec(`UPDATE notification_rules SET priority = $3, provider_id = $4
	WHERE user_id = $1 AND notification_type = $2`, userID, notificationType, priority, providerID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("rows updated", rowsAffected)

	return nil
}

func (r Repository) RemoveNotificationRule(db *sql.DB, id int64) error {
	res, err := db.Exec(`DELETE FROM notification_rules WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("rows deleted", rowsAffected)

	return nil
}

func (r Repository) ListNotificationRuleByUser(db *sql.DB, userID int64) (*NotificationRule, error) {
	query := `SELECT id, user_id, provider_id, coalesce(chat_id, ''), notification_type FROM notification_rules WHERE user_id = $1 ORDER BY priority LIMIT 1`

	var rule NotificationRule

	var chatID string

	err := db.QueryRow(query, userID).Scan(&rule.ID, &rule.UserID, &rule.ProviderID, &chatID, &rule.NotificationType)
	if err != nil {
		return nil, err
	}

	if chatID != "" {
		var err error
		rule.ChatID, err = strconv.ParseInt(chatID, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return &rule, nil
}

func (r Repository) ListNotificationRules(db *sql.DB) ([]*NotificationRule, error) {
	query := fmt.Sprintf(`SELECT id, user_id, notification_type, provider_id, coalesce(chat_id, ''), priority FROM notification_rules`)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []*NotificationRule

	for rows.Next() {
		var (
			id               int64
			userID           int64
			notificationType string
			providerID       string
			chatID           string
			priority         int
		)

		err := rows.Scan(&id, &userID, &notificationType, &providerID, &chatID, &priority)
		if err != nil {
			return nil, err
		}

		var chatIDVal int64
		if chatID != "" {
			chatIDVal, _ = strconv.ParseInt(chatID, 10, 64)
		}

		result = append(result, &NotificationRule{
			ID:               id,
			UserID:           userID,
			NotificationType: notificationType,
			ProviderID:       providerID,
			ChatID:           chatIDVal,
			Priority:         priority,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
