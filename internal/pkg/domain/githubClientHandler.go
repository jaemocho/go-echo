package domain

import (
	"context"
	"log"

	"github.com/google/go-github/v50/github"
)

type GithubClientHandler struct {
	Client *github.Client
}

func NewGithubClientHandler() *GithubClientHandler {
	client := github.NewClient(nil)

	return &GithubClientHandler{
		Client: client,
	}
}

func (g *GithubClientHandler) GetRepoList(owner string) []*github.Repository {
	repos, _, err := g.Client.Repositories.List(context.Background(), owner, nil)
	if _, ok := err.(*github.RateLimitError); ok {
		log.Println("hit rate limit")
		return nil
	}
	return repos
}

func (g *GithubClientHandler) GetWorkflowList(owner, repo string) []*github.Workflow {
	opts := &github.ListOptions{Page: 1, PerPage: 1}
	workflows, _, err := g.Client.Actions.ListWorkflows(context.Background(), owner, repo, opts)
	if err != nil {
		log.Printf("Actions.ListWorkflows returned error: %v", err)
		return nil
	}

	if i := workflows.TotalCount; *i > 0 {
		return workflows.Workflows
	}
	return nil
}
