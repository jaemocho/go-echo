package domain

import (
	"backend/config"
	"context"
	"log"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type GithubClientHandler struct {
	Client *github.Client
}

func NewGithubClientHandler(cfg config.Config) *GithubClientHandler {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &GithubClientHandler{
		Client: client,
	}
}

// onwer를 넣으면 url 기반으로 가져와서 private 이 보이지 않고
//
// owner를 넣지 않으면 token 기반으로 가져와서 private 까지 확인 가능
func (g *GithubClientHandler) GetRepoList(owner string) ([]*github.Repository, error) {
	opts := &github.RepositoryListOptions{Type: "owner", Sort: "updated", Direction: "desc"}
	repos, _, err := g.Client.Repositories.List(context.Background(), owner, opts)
	if err != nil {
		log.Printf("Repositories.List returned error: %v", err)
		return nil, err
	}
	return repos, err
}

func (g *GithubClientHandler) GetWorkflowList(owner, repo string) ([]*github.Workflow, error) {

	workflows, _, err := g.Client.Actions.ListWorkflows(context.Background(), owner, repo, nil)
	if err != nil {
		log.Printf("Actions.ListWorkflows returned error: %v", err)
		return nil, err
	}

	if i := workflows.TotalCount; *i > 0 {
		return workflows.Workflows, err
	}
	return nil, err
}

func (g *GithubClientHandler) CreateWorkflowDispatchEventByFileName(owner, repo, workflowFileName, branch string, inputs map[string]interface{}) error {

	event := github.CreateWorkflowDispatchEventRequest{
		Ref:    branch,
		Inputs: inputs,
		//  Inputs: map[string]interface{}{
		//  	"key": "value",
		//  },
	}

	_, err := g.Client.Actions.CreateWorkflowDispatchEventByFileName(context.Background(), owner, repo, workflowFileName, event)
	if err != nil {
		log.Printf("Actions.CreateWorkflowDispatchEventByFileName returned error: %v", err)
		return err
	}

	return nil
}

func (g *GithubClientHandler) CreateRepo(name, description string, isPrivate, isAutoInt bool) (*github.Repository, error) {

	r := &github.Repository{Name: &name, Private: &isPrivate, Description: &description, AutoInit: &isAutoInt}

	repo, _, err := g.Client.Repositories.Create(context.Background(), "", r)

	if err != nil {
		log.Printf("Repositories.Create returned error: %v", err)
		return nil, err
	}

	return repo, err
}

func (g *GithubClientHandler) DeleteRepo(owner, repo string) error {

	_, err := g.Client.Repositories.Delete(context.Background(), owner, repo)

	if err != nil {
		log.Printf("Repositories.Delete returned error: %v", err)
		return err
	}

	return nil
}
