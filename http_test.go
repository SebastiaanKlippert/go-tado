package tado

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/SebastiaanKlippert/go-tado/tadoauth"
	"github.com/stretchr/testify/assert"
)

type testGet struct{}

func (tg *testGet) method() string {
	return http.MethodGet
}

func (tg *testGet) path() string {
	return "/v1/somepath"
}

func (tg *testGet) body() interface{} {
	return nil
}

type testPost struct {
	InStr string
	InInt int
}

func (tp *testPost) method() string {
	return http.MethodPost
}

func (tp *testPost) path() string {
	return "/v1/somepath"
}

func (tp *testPost) body() interface{} {
	return tp
}

type testOut struct {
	OutField string
}

func TestClient_Do(t *testing.T) {
	// Create mock authentication struct
	tr := &tadoauth.TokenResponse{
		AccessToken:  "thisIsAFakeToken",
		TokenType:    "bearer",
		RefreshToken: "refreshToken",
		ExpiresIn:    588,
		Scope:        "scope",
		Jti:          "jti",
	}
	// Create mock Tado client
	c := NewClient("", "")
	c.tr = tr
	c.accessTokenValidUntil = time.Now().Add(time.Minute)

	// Prepare fake output
	out := new(testOut)

	// Start fake HTTP server
	mockStatusCode := http.StatusBadRequest
	mockResponseBody := "Bad request"
	var incomingRequest *http.Request
	var incomingRequestBody []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		incomingRequest = r
		incomingRequestBody, _ = ioutil.ReadAll(r.Body)
		w.WriteHeader(mockStatusCode)
		_, _ = fmt.Fprint(w, mockResponseBody)
	}))
	defer ts.Close()

	c.baseURL = ts.URL

	// Test a GET method returning a Bad Request status
	tg := new(testGet)

	err := c.do(tg, out)

	if incomingRequest == nil {
		t.Fatal("incomingRequest is nil")
	}

	assert.Equal(t, http.MethodGet, incomingRequest.Method)
	assert.Equal(t, "/v1/somepath", incomingRequest.URL.String())
	assert.Equal(t, "", string(incomingRequestBody))
	assert.Equal(t, "Bearer thisIsAFakeToken", incomingRequest.Header.Get("Authorization"))
	if assert.Error(t, err) {
		assert.Equal(t, "error: HTTP status 400: Bad request", err.Error())
	}

	// Test a GET method returning an OK status with invalid JSON
	mockStatusCode = http.StatusOK
	mockResponseBody = `{ThisIsNotJSON}`

	err = c.do(tg, out)

	if assert.Error(t, err) {
		assert.True(t, strings.HasPrefix(err.Error(), "error decoding output: "), "error has unexpected prefix, error is: %s", err)
	}

	// Test a GET method returning an OK status with valid JSON
	mockStatusCode = http.StatusOK
	mockResponseBody = `{"OutField":"AbCdEf"}`

	err = c.do(tg, out)

	assert.NoError(t, err)
	assert.Equal(t, "AbCdEf", out.OutField)

	// Test a POST method with a valid JSON response
	mockResponseBody = `{"OutField":"AbCdEfGh"}`
	tp := &testPost{
		InStr: "string",
		InInt: 999,
	}

	err = c.do(tp, out)

	assert.NoError(t, err)
	assert.Equal(t, `{"InStr":"string","InInt":999}`+"\n", string(incomingRequestBody))
	assert.Equal(t, http.MethodPost, incomingRequest.Method)
	assert.Equal(t, "application/json", incomingRequest.Header.Get("Content-Type"))
	assert.Equal(t, "AbCdEfGh", out.OutField)
}

func TestInputMethods(t *testing.T) {
	//test that body is not nil for PUT and POST methods and is nil for all other methods
	testStructs := []input{
		new(GetMeInput),
		new(GetHomeInput),
		new(GetHomeStateInput),
		new(GetZonesInput),
		new(GetZoneStateInput),
		new(GetWeatherInput),
		new(PutOverlayInput),
	}

	for _, s := range testStructs {
		switch s.method() {
		case http.MethodPut, http.MethodPost:
			assert.NotNil(t, s.body(), "body is nil on %T with method %s", s, s.method())
		default:
			assert.Nil(t, s.body(), "body is not nil on %T with method %s", s, s.method())
		}
	}
}
