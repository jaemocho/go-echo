package domain

type GitRepo struct {
	Name string `json:"name"`
}

type GitWorkflow struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
}
