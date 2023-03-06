package security

import (
	"backend/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// JwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	// ID   string `json:"id"`
	// Name string `json:"name"`
	jwt.RegisteredClaims
}

type Token struct {
	Token string `json:"token"`
}

var whiteListPaths = []string{
	"/favicon.ico",
	"/swagger/*",
	"/api/v1/login",
	// "/api/*",
	// "/api/v1/signup",
}

func init() {
	middleware.ErrJWTMissing.Code = 401
	middleware.ErrJWTMissing.Message = "Unauthorized"
}

func skipAuth(c echo.Context) bool {
	// Skip authentication for and signup login requests
	for _, path := range whiteListPaths {
		if path == c.Request().URL.Path || path == c.Path() {
			return true
		}
	}
	return false
}

func WebSecurityConfig(e *echo.Echo, cfg config.Config) {

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(cfg.JWTSigningKey),
		Skipper:    skipAuth,
	}
	e.Use(echojwt.WithConfig(config))
}

func JsonWebTokenIssuer(cfg config.Config) string {
	// Set custom claims
	claims := &jwtCustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(cfg.JWTSigningKey))
	if err != nil {
		log.Warn(err)
	}

	return t

}
