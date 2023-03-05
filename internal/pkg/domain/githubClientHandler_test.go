package domain

import (
	"testing"
)

func TestGitHhubHandler(t *testing.T) {
	// assert := assert.New(t)

	gh := NewGithubClientHandler()

	repos := gh.GetRepoList("jaemocho")
	for i, v := range repos {
		t.Log(i, v.GetName())
	}

	workflows := gh.GetWorkflowList("jaemocho", "Study-WebFlux_3")
	for i, v := range workflows {
		t.Log(i, v.GetID(), v.GetName())
	}

}
