package domain

import (
	"backend/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHhubClientHandler(t *testing.T) {
	assert := assert.New(t)

	// token set up
	cfg := config.Config{
		GitHubToken: "",
	}
	gh := NewGithubClientHandler(cfg)

	// repo list 조회 test
	repos, err := gh.GetRepoList("jaemocho")
	assert.NoError(err)
	for i, v := range repos {
		t.Log(i, v.GetName())
		assert.NotNil(v.GetName())
	}

	// workflow list 조회 test
	workflows, err := gh.GetWorkflowList("jaemocho", "Study-WebFlux_3")
	assert.NoError(err)
	for i, v := range workflows {
		t.Log(i, v.GetID(), v.GetName())
		assert.NotNil(v.GetName())
	}

}

func TestGitHhubClientHandler2(t *testing.T) {

	assert := assert.New(t)

	// token set up
	cfg := config.Config{
		GitHubToken: "",
	}
	gh := NewGithubClientHandler(cfg)

	// create repo test
	var (
		name        = "maketest123"
		description = "description"
		private     = false
		autoInit    = false
	)
	repo, err := gh.CreateRepo(name, description, private, autoInit)
	assert.NoError(err)
	assert.Equal("maketest123", *repo.Name)

	// delete repo test
	err = gh.DeleteRepo("jaemocho", "maketest123")
	assert.NoError(err)
}

func TestWorkflowRun(t *testing.T) {
	assert := assert.New(t)

	// token set up
	cfg := config.Config{
		GitHubToken: "",
	}
	gh := NewGithubClientHandler(cfg)

	err := gh.CreateWorkflowDispatchEventByFileName("jaemocho", "Study-WebFlux_3", "maven-publish.yml", "master", nil)
	assert.NoError(err)
}
