package http

import (
	"backend/config"
	"backend/internal/pkg/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GithubHandler struct {
	client *domain.GithubClientHandler
}

func NewGithubHandler(echo *echo.Echo, cfg config.Config) *GithubHandler {

	handler := &GithubHandler{
		client: domain.NewGithubClientHandler(cfg),
	}

	github := echo.Group("/api/v1/github")
	{
		github.GET("/:owner", handler.getReposByOwner)
		github.POST("/:owner", handler.createRepo)
		github.GET("/:owner/:repo", handler.getWorkflowsByRepo)
		github.DELETE("/:owner/:repo", handler.deleteRepo)
	}

	return handler
}

// @Summary		Get repos
// @Description	Get repos by owner
// @name		getReposByOwner
// @Accept		json
// @Produce		json
// @Param		owner	path	string	true	"owner of the repos"
// @Success		200		{array}	domain.GitRepo
// @Router		/api/v1/github/{owner} [get]
// @Security    ApiKeyAuth
func (g *GithubHandler) getReposByOwner(c echo.Context) error {

	owner := c.Param("owner")

	repos, err := g.client.GetRepoList(owner)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	gitRepo := make([]domain.GitRepo, len(repos))

	for i, v := range repos {
		gitRepo[i].Name = v.GetName()
	}

	return c.JSON(http.StatusOK, gitRepo)
}

// @Summary		Get workflows
// @Description	Get workflows by owner, repo
// @name		getWorkflowsByRepo
// @Accept		json
// @Produce		json
// @Param		owner	path	string	true	"owner of the repo"
// @Param		repo	path	string	true	"repo of the workflows"
// @Success		200		{array}	domain.GitWorkflow
// @Router		/api/v1/github/{owner}/{repo} [get]
// @Security    ApiKeyAuth
func (g *GithubHandler) getWorkflowsByRepo(c echo.Context) error {

	owner := c.Param("owner")
	repo := c.Param("repo")

	workflows, err := g.client.GetWorkflowList(owner, repo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	gitWorkflow := make([]domain.GitWorkflow, len(workflows))

	for i, v := range workflows {
		gitWorkflow[i].Name = v.GetName()
		gitWorkflow[i].Id = v.GetID()
	}

	return c.JSON(http.StatusOK, gitWorkflow)
}

// @Summary		Create git Repo
// @Description	Create git Repo
// @name		createRepo
// @Accept		json
// @Produce		json
// @Param		owner	path	string	true	"owner of the repo"
// @Param		repo	body	domain.CreateGitRepo	true	"Repo Info body"
// @Success		201		{object} string
// @Router		/api/v1/github/{owner} [post]
// @Security	ApiKeyAuth
func (g *GithubHandler) createRepo(c echo.Context) error {

	gitRepo := new(domain.CreateGitRepo)

	if err := c.Bind(gitRepo); err != nil {
		c.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	repo, err := g.client.CreateRepo(gitRepo.Name, gitRepo.Description, gitRepo.IsPrivate, gitRepo.IsAutoInt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, repo.GetName()+" create success ")

}

// @Summary		Delete git Repo
// @Description	Delete git Repo
// @name		deleteRepo
// @Accept		json
// @Produce		json
// @Param		owner	path	string	true	"owner of the repo"
// @Param		repo	path	string	true	"repo"
// @Success		200		{object}	string
// @Router		/api/v1/github/{owner}/{repo} [delete]
// @Security    ApiKeyAuth
func (g *GithubHandler) deleteRepo(c echo.Context) error {

	owner := c.Param("owner")
	repo := c.Param("repo")

	err := g.client.DeleteRepo(owner, repo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, repo+" delete success")

}
