package http

import (
	"backend/config"
	"backend/internal/pkg/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GitHandler struct {
	client domain.GitClientHandler
}

func NewGitHandler(echo *echo.Echo, cfg config.Config) *GitHandler {

	handler := &GitHandler{
		client: domain.NewGitClientHandler(cfg),
	}

	gitClient := echo.Group("/api/v1/github")
	{
		gitClient.GET("/:owner", handler.getReposByOwner)
		gitClient.POST("/:owner", handler.createRepo)
		gitClient.GET("/:owner/:repo", handler.getWorkflowsByRepo)
		gitClient.DELETE("/:owner/:repo", handler.deleteRepo)

		gitClient.POST("/issue/:owner/:repo", handler.createIssue)
		gitClient.GET("/issue/:owner/:repo", handler.getIssuesByRepo)
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
func (g *GitHandler) getReposByOwner(c echo.Context) error {

	owner := c.Param("owner")

	repos, err := g.client.GetRepoList(owner)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, repos)
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
func (g *GitHandler) getWorkflowsByRepo(c echo.Context) error {

	owner := c.Param("owner")
	repo := c.Param("repo")

	workflows, err := g.client.GetWorkflowList(owner, repo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, workflows)
}

// @Summary		Create git Repo
// @Description	Create git Repo
// @name		createRepo
// @Accept		json
// @Produce		json
// @Param		owner	path	string	true	"owner of the repo"
// @Param		repo	body	domain.CreateGitRepoRequest	true	"Repo Info body"
// @Success		201		{object} string
// @Router		/api/v1/github/{owner} [post]
// @Security	ApiKeyAuth
func (g *GitHandler) createRepo(c echo.Context) error {

	createGitRepoRequest := new(domain.CreateGitRepoRequest)

	if err := c.Bind(createGitRepoRequest); err != nil {
		c.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	repo, err := g.client.CreateRepo(createGitRepoRequest)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, repo.Name+" create success ")

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
func (g *GitHandler) deleteRepo(c echo.Context) error {

	owner := c.Param("owner")
	repo := c.Param("repo")

	err := g.client.DeleteRepo(owner, repo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, repo+" delete success")

}

// @Summary		Create git Repo Issue
// @Description	Create git Repo Issue
// @name		createIssue
// @Accept		json
// @Produce		json
// @Param		owner	path	string	true	"owner of the repo"
// @Param		repo	path	string	true	"repo"
// @Param		issue	body	domain.CreateGitIssueRequest	true	"Issue Info body"
// @Success		201		{object} string
// @Router		/api/v1/github/issue/{owner}/{repo} [post]
// @Security	ApiKeyAuth
func (g *GitHandler) createIssue(c echo.Context) error {

	gitIssue := new(domain.CreateGitIssueRequest)

	if err := c.Bind(gitIssue); err != nil {
		c.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	owner := c.Param("owner")
	repo := c.Param("repo")

	newIssue, err := g.client.CreateIssue(owner, repo, gitIssue)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, newIssue.Title+" create success ")

}

// @Summary		Get Issues by repo
// @Description	Get Issues by repo
// @name		getIssuesByRepo
// @Accept		json
// @Produce		json
// @Param		owner	path	string	true	"owner of the repos"
// @Param		repo	path	string	true	"repo"
// @Success		200		{array}	domain.GitIssue
// @Router		/api/v1/github/issue/{owner}/{repo} [get]
// @Security    ApiKeyAuth
func (g *GitHandler) getIssuesByRepo(c echo.Context) error {

	owner := c.Param("owner")
	repo := c.Param("repo")

	issues, err := g.client.GetIssueList(owner, repo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, issues)
}
