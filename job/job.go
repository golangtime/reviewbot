package job

import (
	"database/sql"
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

func (job *Job) groupPullRequestByRepository(repo string) {

}

func (job *Job) Run() error {
	job.logger.Info("run job")
	g := job.git

	// TODO: for test (not actually needed)
	// _, err := g.ListRepositories(repoOwner)
	// if err != nil {
	// 	job.logger.Error("list repositories error", "error", err)
	// }

	// fetch all repositories and owners
	repos, err := job.repo.ListRepos(job.db, "")
	if err != nil {
		return err
	}

	if len(repos) == 0 {
		job.logger.Warn("no repositories")
	}

	for _, r := range repos {
		_, err = g.ListPullRequests(r.Owner, r.Name)
		if err != nil {
			job.logger.Error("list pull request error", "error", err)
		}
	}

	// for _, repo := range repos {
	// 	mergeRequests := g.GetOpenMergeRequests()
	// 	for _, request := range mergeRequests {
	// 		peers, err := job.CheckMergeRequest(request)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		for _, peer := range peers {
	// 			job.logger.Info("peer must review open merge request", "peer", peer)
	// 		}
	// 	}
	// }

	return nil
}

// CheckMergeRequest проверяет нужно ли оповещать участников по заданному MergeRequest
func (b *Job) CheckMergeRequest() ([]string, error) {
	// определить минимальное количество апрувов

	// определить состав участников

	// подсчитать сколько участников поставили approve

	// если количество меньше заданного то добавить людей в очередь оповещений

	return nil, nil
}

// PrepareNotifications определяет мин
func (b *Job) PrepareNotifications() {

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
