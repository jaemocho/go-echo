package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v50/github"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetRepos(t *testing.T) {

	e := echo.New()

	gh := NewGithubHandler(e)

	// 1. 조회 테스트 reop
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/:owner")
	c.SetParamNames("owner")
	c.SetParamValues("jaemocho")

	if assert.NoError(t, gh.getReposByOwner(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// rec 에서 읽어올 struct 생성
	repos := []*github.Repository{}

	// json decoder를 이용하여 decoding 반환 값이 domain user의 slice형태
	err := json.NewDecoder(rec.Body).Decode(&repos)
	assert.NoError(t, err)

	for _, v := range repos {
		t.Log(v.GetID(), v.GetName())
	}

}

func TestGetWorkflows(t *testing.T) {

	e := echo.New()

	gh := NewGithubHandler(e)

	// 1. 조회 테스트 reop
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/:owner/:repo")
	c.SetParamNames("owner", "repo")
	c.SetParamValues("jaemocho", "Study-WebFlux_3")

	if assert.NoError(t, gh.getWorkflowsByRepo(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// rec 에서 읽어올 struct 생성
	workflows := []*github.Workflow{}

	// json decoder를 이용하여 decoding 반환 값이 domain user의 slice형태
	err := json.NewDecoder(rec.Body).Decode(&workflows)
	assert.NoError(t, err)

	for _, v := range workflows {
		t.Log(v.GetID(), v.GetName())
	}

}
