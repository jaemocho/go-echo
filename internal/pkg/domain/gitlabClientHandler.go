package domain

import (
	"backend/config"
	"log"
	"strconv"

	"github.com/xanzy/go-gitlab"
)

type GitlabClientHandler struct {
	client *gitlab.Client
}

func NewGitlabClientHandler(cfg config.Config) GitClientHandler {

	client, err := gitlab.NewClient(cfg.GitLabToken)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return &GitlabClientHandler{
		client: client,
	}
}

func (g *GitlabClientHandler) GetRepoList(owner string) ([]*GitRepo, error) {

	projects, _, err := g.client.Projects.ListUserProjects(owner, nil)
	if err != nil {
		log.Printf("Projects.ListProjects returned error: %v", err)
		return nil, err
	}

	gitRepos := make([]*GitRepo, len(projects))

	for i, v := range projects {
		gitRepos[i] = &GitRepo{
			Name:        v.Name,
			Description: v.Description,
			Id:          strconv.Itoa(v.ID),
		}
	}

	return gitRepos, nil
}

func (g *GitlabClientHandler) GetWorkflowList(owner, repo string) ([]*GitWorkflow, error) {

	return nil, nil
}

func (g *GitlabClientHandler) CreateWorkflowDispatchEventByFileName(owner, repo, workflowFileName, branch string, inputs map[string]interface{}) error {

	return nil
}

func (g *GitlabClientHandler) CreateRepo(createGitRepoRequest *CreateGitRepoRequest) (*GitRepo, error) {

	var visibility *gitlab.VisibilityValue
	if !createGitRepoRequest.IsPrivate {
		visibility = gitlab.Visibility(gitlab.PublicVisibility)
	} else {
		visibility = gitlab.Visibility(gitlab.PrivateVisibility)
	}

	opt := &gitlab.CreateProjectOptions{
		Name:                 &createGitRepoRequest.Name,
		Description:          &createGitRepoRequest.Description,
		Visibility:           visibility,
		InitializeWithReadme: &createGitRepoRequest.IsAutoInt,
	}

	project, _, err := g.client.Projects.CreateProject(opt)

	if err != nil {
		log.Printf("Projects.CreateProject returned error: %v", err)
		return nil, err
	}

	gitRepo := &GitRepo{
		Name:        project.Name,
		Description: project.Description,
		IsPrivate:   createGitRepoRequest.IsPrivate,
	}

	return gitRepo, err
}

func (g *GitlabClientHandler) DeleteRepo(owner, repo string) error {

	_, err := g.client.Projects.DeleteProject(owner + "/" + repo)
	if err != nil {
		log.Printf("Repositories.Delete returned error: %v", err)
		return err
	}
	return nil

}

func (g *GitlabClientHandler) CreateIssue(owner, repo string, issueRequest *CreateGitIssueRequest) (*GitIssue, error) {

	var labels gitlab.Labels = issueRequest.Labels

	issue := &gitlab.CreateIssueOptions{
		Title:       &issueRequest.Title,
		Description: &issueRequest.Body,
		Labels:      &labels,
	}

	newIssue, _, err := g.client.Issues.CreateIssue(owner+"/"+repo, issue)

	if err != nil {
		log.Printf("Issues.CreateIssue returned error: %v", err)
		return nil, err
	}

	gitIssue := &GitIssue{
		Title:  newIssue.Title,
		Body:   newIssue.Description,
		Labels: newIssue.Labels,
		Owner:  owner,
		Repo:   repo,
	}

	return gitIssue, err
}

func (g *GitlabClientHandler) GetIssueList(owner, repo string) ([]*GitIssue, error) {

	issueList, _, err := g.client.Issues.ListProjectIssues(owner+"/"+repo, nil)

	if err != nil {
		log.Printf("Issues.ListProjectIssues returned error: %v", err)
		return nil, err
	}

	gitIssueList := make([]*GitIssue, len(issueList))

	for i, v := range issueList {
		gitIssueList[i] = &GitIssue{
			Title:  v.Title,
			Body:   v.Description,
			Labels: v.Labels,
			Owner:  owner,
			Repo:   repo,
		}
	}

	return gitIssueList, nil

}
