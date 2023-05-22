package domain

import "backend/config"

// Github/GitLab/bitbucket client interface
type GitClientHandler interface {
	GetRepoList(owner string) ([]*GitRepo, error)
	GetWorkflowList(owner, repo string) ([]*GitWorkflow, error)
	CreateWorkflowDispatchEventByFileName(owner, repo, workflowFileName, branch string, inputs map[string]interface{}) error
	CreateRepo(createGitRepoRequest *CreateGitRepoRequest) (*GitRepo, error)
	DeleteRepo(owner, repo string) error
	CreateIssue(owner, repo string, issueRequest *CreateGitIssueRequest) (*GitIssue, error)
	GetIssueList(owner, repo string) ([]*GitIssue, error)
}

type GitRepo struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"isPrivate"`
	Id          string `json:"id,omitempty"`
}

type GitWorkflow struct {
	Name string `json:"name,omitempty"`
	Id   int64  `json:"id,omitempty"`
}

type GitIssue struct {
	Title    string   `json:"title,omitempty"`
	Body     string   `json:"body,omitempty"`
	Labels   []string `json:"labels,omitempty"`
	Assignee string   `json:"assignee,omitempty"`
	Owner    string   `json:"owner,omitempty"`
	Repo     string   `json:"repo,omitempty"`
}

type CreateGitRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"isPrivate"`
	IsAutoInt   bool   `json:"isAutoInt"`
}

type CreateGitIssueRequest struct {
	Title    string   `json:"title,omitempty"`
	Body     string   `json:"body,omitempty"`
	Labels   []string `json:"labels,omitempty"`
	Assignee string   `json:"assignee,omitempty"`
}

func NewGitClientHandler(cfg config.Config) GitClientHandler {
	if cfg.GitClient == "github" {
		return NewGithubClientHandler(cfg)
	} else if cfg.GitClient == "gitlab" {
		return NewGitlabClientHandler(cfg)
	} else {
		return NewGithubClientHandler(cfg)
	}

}
