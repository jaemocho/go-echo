package domain

import (
	"backend/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// git hub token setup
	cfg = config.Config{
		GitHubToken: "",
	}
)

func TestGitHhubClientHandler(t *testing.T) {
	assert := assert.New(t)

	gh := NewGithubClientHandler(cfg)

	// repo list 조회 test
	repos, err := gh.GetRepoList("jaemocho")
	assert.NoError(err)
	for i, v := range repos {
		t.Log(i, v.Name)
		assert.NotNil(v.Name)
	}

	// workflow list 조회 test
	workflows, err := gh.GetWorkflowList("jaemocho", "Study-WebFlux_3")
	assert.NoError(err)
	for i, v := range workflows {
		t.Log(i, v.Id, v.Name)
		assert.NotNil(v.Name)
	}

}

func TestGitHhubClientHandler2(t *testing.T) {

	assert := assert.New(t)

	gh := NewGithubClientHandler(cfg)

	createGitRepoRequest := &CreateGitRepoRequest{
		Name:        "maketest123",
		Description: "description",
		IsPrivate:   false,
		IsAutoInt:   false,
	}

	repo, err := gh.CreateRepo(createGitRepoRequest)
	assert.NoError(err)
	assert.Equal("maketest123", repo.Name)

	// delete repo test
	// err = gh.DeleteRepo("jaemocho", "maketest123")
	// assert.NoError(err)
}

func TestWorkflowRun(t *testing.T) {
	assert := assert.New(t)

	gh := NewGithubClientHandler(cfg)

	err := gh.CreateWorkflowDispatchEventByFileName("jaemocho", "Study-WebFlux_3", "maven-publish.yml", "master", nil)
	assert.NoError(err)
}

func TestIssue(t *testing.T) {
	assert := assert.New(t)

	gh := NewGithubClientHandler(cfg)

	title := "test"
	body := "test body"
	assignee := "jaemocho"
	labels := []string{"bug", "number"}
	inputIssue := &CreateGitIssueRequest{Title: title, Body: body, Assignee: assignee, Labels: labels}

	newIssue, err := gh.CreateIssue("jaemocho", "Study-WebFlux_3", inputIssue)
	assert.NoError(err)
	assert.Equal("test", newIssue.Title)

	issueList, err := gh.GetIssueList("jaemocho", "Study-WebFlux_3")
	assert.NoError(err)
	for _, v := range issueList {
		t.Log(v.Title, v.Body, v.Labels)
	}

}
