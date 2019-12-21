package tado

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SebastiaanKlippert/go-tado/tadoauth"
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

func setupTestClientAndServer(hf http.HandlerFunc) (*Client, *httptest.Server) {
	s := httptest.NewServer(hf)
	c := NewClient("username", "password")
	c.accessTokenValidUntil = time.Now().Add(time.Hour) // to ensure we don't go to Tado authentication
	c.baseURL = s.URL
	c.tr = &tadoauth.TokenResponse{
		AccessToken: "fakeToken",
	}
	return c, s
}

func TestClient_GetMe(t *testing.T) {

	called := false
	f := func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "/v1/me", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprintln(w, `{"name":"SK"}`)
	}

	client, server := setupTestClientAndServer(f)
	defer server.Close()

	m, err := client.GetMe()
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
	if assert.NotNil(t, m) {
		assert.Equal(t, "SK", m.Name)
	}
}

func TestClient_GetHome(t *testing.T) {

	called := false
	f := func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "/v2/homes/12345", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprintln(w, `{"id": 12345}`)
	}

	client, server := setupTestClientAndServer(f)
	defer server.Close()

	in := &GetHomeInput{
		HomeID: 12345,
	}

	h, err := client.GetHome(in)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
	if assert.NotNil(t, h) {
		assert.Equal(t, 12345, h.ID)
	}
}

func TestClient_GetZones(t *testing.T) {

	called := false
	f := func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "/v2/homes/12345/zones", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprintln(w, `[{"id": 1}, {"id": 2}]`)
	}

	client, server := setupTestClientAndServer(f)
	defer server.Close()

	in := &GetZonesInput{
		HomeID: 12345,
	}

	z, err := client.GetZones(in)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
	if assert.NotEmpty(t, z) && assert.Equal(t, 2, len(z)) {
		assert.Equal(t, 1, z[0].ID)
		assert.Equal(t, 2, z[1].ID)
	}
}
