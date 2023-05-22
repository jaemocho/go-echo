package domain

import (
	"backend/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// git hub token setup
	cfgLab = config.Config{
		GitLabToken: "",
	}
)

func TestGitlabGetRepos(t *testing.T) {
	assert := assert.New(t)

	gh := NewGitlabClientHandler(cfgLab)

	// repo list 조회 test
	repos, err := gh.GetRepoList("mot882000")
	assert.NoError(err)
	for i, v := range repos {
		t.Log(i, v.Name, v.Description)
		assert.NotNil(v.Name)
	}
}

func TestGitlabCreateRepo(t *testing.T) {

	assert := assert.New(t)

	gh := NewGitlabClientHandler(cfgLab)

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
	err = gh.DeleteRepo("mot882000", "maketest123")
	assert.NoError(err)
}

func TestCreateIssue(t *testing.T) {
	assert := assert.New(t)

	gh := NewGitlabClientHandler(cfgLab)

	title := "test"
	body := "test body"
	assignee := "jaemocho"
	labels := []string{"bug", "number", "ttt"}
	inputIssue := &CreateGitIssueRequest{Title: title, Body: body, Assignee: assignee, Labels: labels}

	newIssue, err := gh.CreateIssue("mot882000", "gitlab-test-project", inputIssue)
	assert.NoError(err)
	assert.Equal("test", newIssue.Title)
	t.Log(newIssue.Body, newIssue.Title, newIssue.Labels, newIssue.Assignee)

	issueList, err := gh.GetIssueList("mot882000", "gitlab-test-project")
	assert.NoError(err)

	for _, v := range issueList {
		t.Log(v.Title, v.Body, v.Labels, v.Owner, v.Repo)
	}

}
