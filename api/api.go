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

func (v *V1) AddRepo(owner, repo string, minApprovals int) error {
	if repo == "" {
		return fmt.Errorf("empty repository name")
	}

	if owner == "" {
		return fmt.Errorf("empty owner name")
	}

	err := v.repo.AddRepo(v.db, owner, repo, minApprovals)
	return err
}

func (v *V1) ListRepo(owner string) ([]db.RepoEntity, error) {
	if owner == "" {
		return nil, fmt.Errorf("empty owner name")
	}

	resp, err := v.repo.ListRepos(v.db, owner)
	return resp, err
}

func (v *V1) DeleteRepo() {

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

func (v *V1) AddNotificationRule(peer string, notifyType string) {

}

func (v *V1) DeactivateNotificationRule(peer string) {

}
