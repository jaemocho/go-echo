package security

import (
	"backend/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonWebTokenIssuer(t *testing.T) {
	assert := assert.New(t)
	cfg := config.Config{
		JWTSigningKey: "signingkey",
	}

	token := JsonWebTokenIssuer(cfg)
	t.Log(token)

	assert.NotNil(token)
}
