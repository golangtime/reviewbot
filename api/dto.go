package api

type AddRepoRequest struct {
	Owner        string `json:"owner"`
	Name         string `json:"name"`
	MinApprovals int    `json:"min_approvals"`
}

type AddRepoResponse struct {
	Success bool `json:"success"`
}

type ListReposRequest struct {
	Owner string `json:"owner"`
}

type ListReposResponse struct {
	Repos []*Repo
}

type Repo struct {
	Name         string `json:"name"`
	MinApprovals int    `json:"min_approvals"`
}
