package tado

import (
	"testing"
	"time"

	"github.com/SebastiaanKlippert/go-tado/tadoauth"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	tr := &tadoauth.TokenResponse{
		AccessToken:  "token",
		TokenType:    "bearer",
		RefreshToken: "refreshToken",
		ExpiresIn:    588,
		Scope:        "scope",
		Jti:          "jti",
	}

	c := NewClient(tr)

	assert.NotNil(t, c.HTTPClient, "HTTP client is nil")
	assert.NotNil(t, c.mutex, "mutex is nil")
	assert.Equal(t, tr, c.tr, "tr is nil")
	assert.Equal(t, defaultBaseURL, c.baseURL, "baseURL is incorrect")
	assert.True(t, time.Now().Before(c.accessTokenValidUntil), "accessTokenValidUntil invalid, have %s", c.accessTokenValidUntil)
}

/*
func TestClient_GetMe(t *testing.T) {
	tr, err := tadoauth.GetToken("", "")
	if err != nil {
		t.Fatal(err)
	}
	c := NewClient(tr)
	m, err := c.GetMe()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%+v\n", m)
}
*/
