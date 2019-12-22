package tado

import (
	"fmt"
	"io/ioutil"
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

type mockAuthClient struct {
	username, password, refreshToken string
}

func (m *mockAuthClient) GetToken(username, password string) (*tadoauth.TokenResponse, error) {
	m.username = username
	m.password = password
	return &tadoauth.TokenResponse{}, nil
}

func (m *mockAuthClient) RefreshToken(refreshToken string) (*tadoauth.TokenResponse, error) {
	m.refreshToken = refreshToken
	return &tadoauth.TokenResponse{}, nil
}

func TestClient_validateAccessToken(t *testing.T) {
	// create new client with expired token
	c := NewClient("username1", "password1")
	c.tr = &tadoauth.TokenResponse{
		RefreshToken: "",
	}

	// construct mock auth client
	mockAuth := new(mockAuthClient)
	c.authClient = mockAuth

	// use testserver that returns an empty json response
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{}`)
	}))
	c.baseURL = s.URL

	// call any method
	_, err := c.GetMe()

	// check if authentication methods have been called
	assert.NoError(t, err)
	assert.Equal(t, "username1", mockAuth.username)
	assert.Equal(t, "password1", mockAuth.password)

	// now check if refreshtoken method is being called
	c.accessTokenValidUntil = time.Time{}
	c.tr = &tadoauth.TokenResponse{
		RefreshToken: "fakeRefreshToken",
	}

	_, err = c.GetMe()

	assert.NoError(t, err)
	assert.Equal(t, "fakeRefreshToken", mockAuth.refreshToken)
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
		assert.Equal(t, "/v2/me", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `{"name":"SK"}`)
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
		_, _ = fmt.Fprint(w, `{"id": 12345}`)
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
		_, _ = fmt.Fprint(w, `[{"id": 1}, {"id": 2}]`)
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

func TestClient_GetHomeState(t *testing.T) {

	called := false
	f := func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "/v2/homes/12345/state", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `{"presence": "HOME"}`)
	}

	client, server := setupTestClientAndServer(f)
	defer server.Close()

	in := &GetHomeStateInput{
		HomeID: 12345,
	}

	s, err := client.GetHomeState(in)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
	if assert.NotNil(t, s) {
		assert.Equal(t, HomeStateHome, s.Presence)
	}
}

func TestClient_GetZoneState(t *testing.T) {

	called := false
	f := func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "/v2/homes/12345/zones/2/state", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `{"tadoMode": "HOME"}`)
	}

	client, server := setupTestClientAndServer(f)
	defer server.Close()

	in := &GetZoneStateInput{
		HomeID: 12345,
		ZoneID: 2,
	}

	s, err := client.GetZoneState(in)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
	if assert.NotNil(t, s) {
		assert.Equal(t, "HOME", s.TadoMode)
	}
}

func TestClient_GetWeather(t *testing.T) {

	called := false
	f := func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "/v2/homes/12345/weather", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `{"outsideTemperature": {"celsius": 8.50}}`)
	}

	client, server := setupTestClientAndServer(f)
	defer server.Close()

	in := &GetWeatherInput{
		HomeID: 12345,
	}

	w, err := client.GetWeather(in)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
	if assert.NotNil(t, w) {
		assert.Equal(t, 8.5, w.OutsideTemperature.Celsius)
	}
}

func TestClient_GetDayReport(t *testing.T) {

	called := false
	f := func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "/v2/homes/12345/zones/2/dayReport", r.URL.Path)
		assert.Equal(t, "date=2020-12-31", r.URL.RawQuery)
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `{"zoneType": "HEATING"}`)
	}

	client, server := setupTestClientAndServer(f)
	defer server.Close()

	in := &GetDayReportInput{
		HomeID: 12345,
		ZoneID: 2,
		Date:   time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC),
	}

	r, err := client.GetDayReport(in)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
	if assert.NotNil(t, r) {
		assert.Equal(t, "HEATING", r.ZoneType)
	}
}

func TestClient_PutOverlay(t *testing.T) {

	called := false
	f := func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "/v2/homes/12345/zones/99/overlay", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		b, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, `{"setting":{"type":"HEATING","power":"ON","temperature":{"celsius":17}},"termination":{"type":"MANUAL"}}`+"\n", string(b))
		_, _ = fmt.Fprint(w, `{"type": "MANUAL"}`)
	}

	client, server := setupTestClientAndServer(f)
	defer server.Close()

	in := &PutOverlayInput{
		HomeID: 12345,
		ZoneID: 99,
		OverlayInput: OverlayInput{
			Setting: OverlayInputSetting{
				Type:  "HEATING",
				Power: "ON",
				Temperature: OverlayInputTemperature{
					Celsius: 17,
				},
			},
			Termination: OverlayInputTermination{
				Type: TerminationTypeManual,
			},
		},
	}

	o, err := client.PutOverlay(in)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
	if assert.NotNil(t, o) {
		assert.Equal(t, "MANUAL", o.Type)
	}
}

func TestClient_DeleteOverlay(t *testing.T) {

	called := false
	f := func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "/v2/homes/12345/zones/2/overlay", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	}

	client, server := setupTestClientAndServer(f)
	defer server.Close()

	in := &DeleteOverlayInput{
		HomeID: 12345,
		ZoneID: 2,
	}

	r, err := client.DeleteOverlay(in)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
	if assert.NotNil(t, r) {
		assert.Empty(t, r)
	}
}
