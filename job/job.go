package job

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/golangtime/reviewbot/client"
	"github.com/golangtime/reviewbot/db"
)

type Job struct {
	repo   *db.Repository
	db     *sql.DB
	git    client.GitClient
	logger *slog.Logger
}

func (job *Job) Run() error {
	job.logger.Info("run job")
	g := job.git

	// fetch all repositories and owners
	repos, err := job.repo.ListRepos(job.db, "")
	if err != nil {
		return err
	}

	if len(repos) == 0 {
		job.logger.Warn("no repositories")
	}

	for _, r := range repos {
		pullRequests, err := g.ListPullRequests(r.Owner, r.Name)
		if err != nil {
			job.logger.Error("list pull request error", "error", err)
		}

		for _, pr := range pullRequests {
			for _, u := range pr.Reviewers {
				fmt.Println("pending review", u.ID, u.Email)
			}

			reviews, err := g.ListReviews(r.Owner, r.Name, int(pr.ExternalID))
			if err != nil {
				job.logger.Error("list reviews error", "error", err)
			}

			countPending := len(pr.Reviewers)
			for _, r := range reviews {
				if r.Status == "APPROVED" {
					countPending--
				}
			}

			if countPending < r.MinApprovals {
				for _, u := range pr.Reviewers {
					email := u.Email
					job.logger.Info("enqueue notification", "url", pr.Link, "email", email, "user_id", u.ID)
					err = job.repo.EnqueueNotification(job.db, "github", pr.Link, email, u.ID)
					if err != nil {
						job.logger.Error("enqueue notification", "error", err)
					}
				}
			}
		}
	}

	return nil
}

func defaultJob(dbConn *sql.DB, logger *slog.Logger, gitClient client.GitClient) *Job {
	repo := &db.Repository{}

	return &Job{
		repo:   repo,
		db:     dbConn,
		git:    gitClient,
		logger: logger,
	}
}

func NewJob(db *sql.DB, logger *slog.Logger, gitClient client.GitClient) func() {
	jb := defaultJob(db, logger, gitClient)
	return func() {
		err := jb.Run()
		if err != nil {
			logger.Error("job error", "error", err)
		}
	}
}
