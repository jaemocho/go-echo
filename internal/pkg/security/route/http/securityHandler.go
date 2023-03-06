package http

import (
	"backend/config"
	"backend/internal/pkg/security"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SecurityHandler struct {
	cfg config.Config
}

func NewSecurityHandler(echo *echo.Echo, cfg config.Config) *SecurityHandler {

	handler := &SecurityHandler{cfg}

	security := echo.Group("/api/v1/login")
	{
		security.GET("", handler.getAccessToken)
	}

	return handler
}

//	@Summary		login (issue token)
//	@Description	get access token
//	@name			getAccessToken
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string
//	@Router			/api/v1/login [get]
func (s *SecurityHandler) getAccessToken(c echo.Context) error {

	token := security.JsonWebTokenIssuer(s.cfg)

	return c.JSON(http.StatusCreated, token)
}
