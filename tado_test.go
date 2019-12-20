package tado

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c := NewClient("", "")

	assert.NotNil(t, c.HTTPClient, "HTTP client is nil")
	assert.NotNil(t, c.mutex, "mutex is nil")
	assert.Nil(t, c.tr, "tr is not nil")
	assert.Equal(t, c.accessTokenValidUntil, time.Time{})
	assert.Equal(t, defaultBaseURL, c.baseURL, "baseURL is incorrect")
}

/*
func TestClient_GetMe(t *testing.T) {
	c := NewClient("", "")
	m, err := c.GetMe()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%+v\n", m)
}
*/
