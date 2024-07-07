package job

import (
	"database/sql"
	"log/slog"

	"github.com/golangtime/reviewbot/client"
	"github.com/golangtime/reviewbot/db"
)

type GitClients struct {
	Github    client.GitClient
	Bitbucket client.GitClient
}

type Job struct {
	repo   *db.Repository
	db     *sql.DB
	git    GitClients
	logger *slog.Logger
}

func (job *Job) doForRepository(r db.RepoEntity, client client.GitClient) error {
	pullRequests, err := client.UnfinishedPullRequests(r.Owner, r.Name, r.MinApprovals)
	if err != nil {
		return err
	}

	for _, pr := range pullRequests {
		for _, u := range pr.Reviewers {
			email := u.Email
			job.logger.Info("enqueue notification", "url", pr.Link, "email", email, "user_id", u.ID)
			err = job.repo.EnqueueNotification(job.db, "github", pr.Link, email, u.ID)
			if err != nil {
				job.logger.Error("enqueue notification", "error", err)
			}
		}
	}

	return nil
}

func (job *Job) Run() error {
	job.logger.Info("run job")

	repos, err := job.repo.ListRepos(job.db, "")
	if err != nil {
		return err
	}

	if len(repos) == 0 {
		job.logger.Warn("no repositories")
	}

	for _, r := range repos {
		var err error
		switch r.Provider {
		case "github":
			err = job.doForRepository(r, job.git.Github)
		case "bitbucket":
			err = job.doForRepository(r, job.git.Bitbucket)
		default:
			job.logger.Warn("not supported provider", "provider", r.Provider)
		}

		if err != nil {
			job.logger.Error("repo error", "error", err)
		}
	}

	return nil
}

func defaultJob(dbConn *sql.DB, logger *slog.Logger, gitClients GitClients) *Job {
	repo := &db.Repository{}

	return &Job{
		repo:   repo,
		db:     dbConn,
		git:    gitClients,
		logger: logger,
	}
}

func NewJob(db *sql.DB, logger *slog.Logger, gitClients GitClients) func() {
	jb := defaultJob(db, logger, gitClients)
	return func() {
		err := jb.Run()
		if err != nil {
			logger.Error("job error", "error", err)
		}
	}
}
