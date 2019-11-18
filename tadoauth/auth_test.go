package tadoauth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	hf http.HandlerFunc
}

func (th *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	th.hf(w, r)
}

func TestGetToken(t *testing.T) {
	incomingBody := ""
	responseBody := ""
	responseStatusCode := 0

	// create test handler for mock auth server
	th := &testHandler{
		hf: func(w http.ResponseWriter, r *http.Request) {
			s, _ := ioutil.ReadAll(r.Body)
			incomingBody = string(s)
			w.WriteHeader(responseStatusCode)
			_, _ = fmt.Fprint(w, responseBody)
		},
	}

	// start mock server
	testServer := httptest.NewServer(th)
	defer testServer.Close()
	// set mock endpoint
	endpoint = testServer.URL
	defer func() { endpoint = defaultEndpoint }()

	// test incoming request with bad credentials
	username, password := "fake@sample.com", "PassW0rd"
	responseStatusCode = http.StatusBadRequest
	responseBody = `{"error": "invalid_grant", "error_description": "Bad credentials"}`

	tokenResponse, err := GetToken(username, password)

	assert.Equal(t, "client_id=tado-web-app&client_secret=wZaRN7rpjn3FoNyF5IFuxg9uMzYJcvOoQ8QWiIqS3hfk6gLhVlG57j5YNoZL2Rtc&grant_type=password&password=PassW0rd&scope=home.user&username=fake%40sample.com", incomingBody)
	if assert.Error(t, err) {
		assert.Equal(t, "authentication error invalid_grant: Bad credentials", err.Error())
	}
	assert.Nil(t, tokenResponse)

	// test incoming request with non-json response
	responseStatusCode = http.StatusInternalServerError
	responseBody = "Bad Gateway"

	tokenResponse, err = GetToken(username, password)

	if assert.Error(t, err) {
		assert.Equal(t, "authentication error: HTTP status 500: Bad Gateway", err.Error())
	}
	assert.Nil(t, tokenResponse)

	// test against wrong endpoint
	endpoint = "invalidhost"

	tokenResponse, err = GetToken(username, password)

	if assert.Error(t, err) {
		assert.True(t, strings.HasPrefix(err.Error(), "authentication HTTP error: "))
	}
	assert.Nil(t, tokenResponse)

	// test valid response code with invalid json
	endpoint = testServer.URL
	responseStatusCode = http.StatusOK
	responseBody = "Not JSON"

	tokenResponse, err = GetToken(username, password)

	if assert.Error(t, err) {
		assert.True(t, strings.HasPrefix(err.Error(), "authentication JSON error: "), "Prefix does not match, have %s", err.Error())
	}
	assert.Nil(t, tokenResponse)

	// test valid response
	responseStatusCode = http.StatusOK
	responseBody = `{
	"access_token": "ABCdef",
    "token_type": "bearer",
    "refresh_token": "fedCBA",
    "expires_in": 599,
    "scope": "home.user",
    "jti": "0000"
	}`

	tokenResponse, err = GetToken(username, password)
	assert.NoError(t, err)

	if assert.NotNil(t, tokenResponse) {
		assert.Equal(t, "ABCdef", tokenResponse.AccessToken)
		assert.Equal(t, "fedCBA", tokenResponse.RefreshToken)
		assert.Equal(t, "home.user", tokenResponse.Scope)
		assert.Equal(t, 599, tokenResponse.ExpiresIn)
	}
}

func TestRefreshToken(t *testing.T) {
	incomingBody := ""
	responseBody := ""
	responseStatusCode := 0

	// create test handler for mock auth server
	th := &testHandler{
		hf: func(w http.ResponseWriter, r *http.Request) {
			s, _ := ioutil.ReadAll(r.Body)
			incomingBody = string(s)
			w.WriteHeader(responseStatusCode)
			_, _ = fmt.Fprint(w, responseBody)
		},
	}

	// start mock server
	testServer := httptest.NewServer(th)
	endpoint = testServer.URL
	defer func() { endpoint = defaultEndpoint }()

	// test incoming request
	responseStatusCode = http.StatusOK
	responseBody = `{"access_token": "t0ken"}`

	tokenResponse, err := RefreshToken("ABCdef123")

	assert.Equal(t, "client_id=tado-web-app&client_secret=wZaRN7rpjn3FoNyF5IFuxg9uMzYJcvOoQ8QWiIqS3hfk6gLhVlG57j5YNoZL2Rtc&grant_type=refresh_token&refresh_token=ABCdef123&scope=home.user", incomingBody)
	assert.NoError(t, err)
	if assert.NotNil(t, tokenResponse) {
		assert.Equal(t, "t0ken", tokenResponse.AccessToken)
	}
}
