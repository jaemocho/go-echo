package http

import (
	"backend/internal/pkg/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GithubHandler struct {
	client *domain.GithubClientHandler
}

func NewGithubHandler(echo *echo.Echo) *GithubHandler {

	handler := &GithubHandler{
		client: domain.NewGithubClientHandler(),
	}

	user := echo.Group("/api/v1/github")
	{
		user.GET("/:owner", handler.getReposByOwner)
		user.GET("/:owner/:repo", handler.getWorkflowsByRepo)
	}

	return handler
}

// @Summary					Get repos
// @Description				Get repos by owner
// @name						getReposByOwner
// @Accept						json
// @Produce					json
// @Param						owner	path	string	true	"owner of the repos"
// @Success					200		{array}	domain.GitRepo
// @Router						/api/v1/github/{owner} [get]
// @Security    ApiKeyAuth
func (g *GithubHandler) getReposByOwner(c echo.Context) error {

	owner := c.Param("owner")

	repos := g.client.GetRepoList(owner)

	gitRepo := make([]domain.GitRepo, len(repos))

	for i, v := range repos {
		gitRepo[i].Name = v.GetName()
	}

	return c.JSON(http.StatusOK, gitRepo)
}

// @Summary					Get workflows
// @Description				Get workflows by owner, repo
// @name						getWorkflowsByRepo
// @Accept						json
// @Produce					json
// @Param						owner	path	string	true	"owner of the repo"
// @Param						repo	path	string	true	"repo of the workflows"
// @Success					200		{array}	domain.GitWorkflow
// @Router						/api/v1/github/{owner}/{repo} [get]
// @Security    ApiKeyAuth
func (g *GithubHandler) getWorkflowsByRepo(c echo.Context) error {

	owner := c.Param("owner")
	repo := c.Param("repo")

	workflows := g.client.GetWorkflowList(owner, repo)

	gitWorkflow := make([]domain.GitWorkflow, len(workflows))

	for i, v := range workflows {
		gitWorkflow[i].Name = v.GetName()
		gitWorkflow[i].Id = v.GetID()
	}

	return c.JSON(http.StatusOK, gitWorkflow)
}
