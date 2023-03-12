package domain

import (
	"backend/config"
	"context"
	"log"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type GithubClientHandler struct {
	client *github.Client
}

func NewGithubClientHandler(cfg config.Config) GitClientHandler {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &GithubClientHandler{
		client: client,
	}
}

// onwer를 넣으면 url 기반으로 가져와서 private 이 보이지 않고
//
// owner를 넣지 않으면 token 기반으로 가져와서 private 까지 확인 가능
func (g *GithubClientHandler) GetRepoList(owner string) ([]*GitRepo, error) {
	opts := &github.RepositoryListOptions{Type: "owner", Sort: "updated", Direction: "desc"}
	repos, _, err := g.client.Repositories.List(context.Background(), owner, opts)
	if err != nil {
		log.Printf("Repositories.List returned error: %v", err)
		return nil, err
	}

	gitRepo := make([]*GitRepo, len(repos))

	for i, v := range repos {
		gitRepo[i] = &GitRepo{
			Name:        v.GetName(),
			Description: v.GetDescription(),
		}
	}

	return gitRepo, err
}

func (g *GithubClientHandler) GetWorkflowList(owner, repo string) ([]*GitWorkflow, error) {

	workflows, _, err := g.client.Actions.ListWorkflows(context.Background(), owner, repo, nil)
	if err != nil {
		log.Printf("Actions.ListWorkflows returned error: %v", err)
		return nil, err
	}

	if cnt := *workflows.TotalCount; cnt > 0 {
		gitWorkFlow := make([]*GitWorkflow, cnt)
		for i, v := range workflows.Workflows {
			gitWorkFlow[i] = &GitWorkflow{
				Id:   v.GetID(),
				Name: v.GetName(),
			}

		}
		return gitWorkFlow, nil
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

	_, err := g.client.Actions.CreateWorkflowDispatchEventByFileName(context.Background(), owner, repo, workflowFileName, event)
	if err != nil {
		log.Printf("Actions.CreateWorkflowDispatchEventByFileName returned error: %v", err)
		return err
	}

	return nil
}

func (g *GithubClientHandler) CreateRepo(name, description string, isPrivate, isAutoInit bool) (*GitRepo, error) {

	r := &github.Repository{
		Name:        &name,
		Private:     &isPrivate,
		Description: &description,
		AutoInit:    &isAutoInit,
	}

	repo, _, err := g.client.Repositories.Create(context.Background(), "", r)

	if err != nil {
		log.Printf("Repositories.Create returned error: %v", err)
		return nil, err
	}

	gitRepo := &GitRepo{
		Name:        repo.GetName(),
		Description: repo.GetDescription(),
		IsPrivate:   *repo.Private,
	}

	return gitRepo, err
}

func (g *GithubClientHandler) DeleteRepo(owner, repo string) error {

	_, err := g.client.Repositories.Delete(context.Background(), owner, repo)

	if err != nil {
		log.Printf("Repositories.Delete returned error: %v", err)
		return err
	}

	return nil
}

func (g *GithubClientHandler) CreateIssue(owner, repo string, issueRequest *CreateGitIssueRequest) (*GitIssue, error) {

	issue := &github.IssueRequest{
		Title:    &issueRequest.Title,
		Body:     &issueRequest.Body,
		Assignee: &issueRequest.Assignee,
		Labels:   &issueRequest.Labels,
	}

	newIssue, _, err := g.client.Issues.Create(context.Background(), owner, repo, issue)

	if err != nil {
		log.Printf("Issues.Create returned error: %v", err)
		return nil, err
	}

	gitIssue := &GitIssue{
		Title:    *newIssue.Title,
		Body:     *newIssue.Body,
		Labels:   *issue.Labels,
		Assignee: *issue.Assignee,
		Owner:    owner,
		Repo:     repo,
	}

	return gitIssue, err
}

func (g *GithubClientHandler) GetIssueList(owner, repo string) ([]*GitIssue, error) {
	opts := &github.IssueListByRepoOptions{Sort: "created", Direction: "desc"}

	issueList, _, err := g.client.Issues.ListByRepo(context.Background(), owner, repo, opts)

	if err != nil {
		log.Printf("Issues.ListByRepo returned error: %v", err)
		return nil, err
	}

	gitIssueList := make([]*GitIssue, len(issueList))

	for i, v := range issueList {
		gitIssueList[i] = &GitIssue{
			Title:  *v.Title,
			Body:   *v.Body,
			Labels: parseIssueLabels(v.Labels),
			Owner:  owner,
			Repo:   repo,
		}
	}

	return gitIssueList, nil

}

func parseIssueLabels(labels []*github.Label) []string {

	returnLabels := make([]string, len(labels))

	for i, v := range labels {
		returnLabels[i] = *v.Name
	}

	return returnLabels
}
