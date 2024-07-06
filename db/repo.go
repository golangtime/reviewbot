package db

import (
	"database/sql"
)

type RepoEntity struct {
	Owner        string
	Name         string
	MinApprovals int
}

type Repo interface {
	ListRepos(db *sql.DB, owner string) ([]RepoEntity, error)
	AddRepo(db *sql.DB, owner, repo string, minApproval int) error
}

type Repository struct {
}

func (r Repository) ListRepos(db *sql.DB, owner string) ([]RepoEntity, error) {
	query := "SELECT name, owner, min_approvals FROM repositories"
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
			minApprovals int
		)

		err := rows.Scan(&name, &owner, &minApprovals)
		if err != nil {
			return nil, err
		}

		repos = append(repos, RepoEntity{
			Name:         name,
			Owner:        owner,
			MinApprovals: minApprovals,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}

func (r Repository) AddRepo(db *sql.DB, owner, repo string, minApproval int) error {
	_, err := db.Exec("INSERT INTO repositories (owner, name, min_approvals) VALUES ($1, $2, $3)", owner, repo, minApproval)
	return err
}
