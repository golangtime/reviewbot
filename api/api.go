package api

import (
	"database/sql"
	"fmt"

	"github.com/golangtime/reviewbot/db"
)

type V1 struct {
	db   *sql.DB
	repo db.Repo
}

func NewAPIV1(db *sql.DB, repo db.Repo) *V1 {
	return &V1{
		db:   db,
		repo: repo,
	}
}

func (v *V1) AddRepo(owner, repo, provider string, minApprovals int) error {
	if repo == "" {
		return fmt.Errorf("empty repository name")
	}

	if owner == "" {
		return fmt.Errorf("empty owner name")
	}

	err := v.repo.AddRepo(v.db, owner, repo, minApprovals, provider)
	return err
}

func (v *V1) ListRepo(owner string) ([]db.RepoEntity, error) {
	resp, err := v.repo.ListRepos(v.db, owner)
	return resp, err
}

func (v *V1) ListPendingNotifications() ([]db.Notification, error) {
	result, err := v.repo.ListPendingNotifications(v.db)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *V1) RemoveRepo() {

}

func (v *V1) MuteRepo() {

}

func (v *V1) UnmuteRepo() {
}

func (v *V1) MutePeer() {

}

func (v *V1) UnmutePeer() {
}

func (v *V1) AddMergeRequesRule() {

}

func (v *V1) AddNotificationRule(userID int64, notifyType string, providerID, chatID string, priority int) error {
	return v.repo.AddNotificationRule(v.db, userID, notifyType, providerID, chatID, priority)
}

func (v *V1) ListNotificationRules() ([]*db.NotificationRule, error) {
	return v.repo.ListNotificationRules(v.db)
}

func (v *V1) DeactivateNotificationRule(peer string) {

}
