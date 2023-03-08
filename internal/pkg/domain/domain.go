package domain

type GitRepo struct {
	Name string `json:"name"`
}

type GitWorkflow struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
}

type CreateGitRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"isPrivate"`
	IsAutoInt   bool   `json:"isAutoInt"`
}
