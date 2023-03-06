package http

import (
	"backend/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHandler(t *testing.T) {

	e := echo.New()
	cfg := config.Config{
		JWTSigningKey: "jwtkey",
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	sh := NewSecurityHandler(e, cfg)

	// Assertions
	if assert.NoError(t, sh.getAccessToken(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	t.Log(rec.Body.String())
}
