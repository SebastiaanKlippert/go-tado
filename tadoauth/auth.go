package tadoauth

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultEndpoint = "https://auth.tado.com/oauth/token"
	clientID        = "tado-web-app"
	clientSecret    = "wZaRN7rpjn3FoNyF5IFuxg9uMzYJcvOoQ8QWiIqS3hfk6gLhVlG57j5YNoZL2Rtc"
	scope           = "home.user"
)

var endpoint = defaultEndpoint

func do(data url.Values) (*TokenResponse, error) {
	// set common form data
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("scope", scope)

	// post form
	resp, err := http.PostForm(endpoint, data)
	if err != nil {
		return nil, fmt.Errorf("authentication HTTP error: %s", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// check HTTP status
	if resp.StatusCode >= http.StatusBadRequest {
		// not OK, read body
		body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<14))
		if err != nil {
			return nil, fmt.Errorf("authentication error: %s", err)
		}
		// check if it is a JSON error response
		ae := new(AuthenticationError)
		err = json.Unmarshal(body, ae)
		if err == nil {
			//no unmarshal error, this is a JSON error response, return this as error
			return nil, ae
		}
		// return the body as error
		return nil, fmt.Errorf("authentication error: HTTP status %d: %s", resp.StatusCode, string(body))
	}

	// HTTP status is OK or similar
	tr := new(TokenResponse)
	err = json.NewDecoder(resp.Body).Decode(tr)
	if err != nil {
		return nil, fmt.Errorf("authentication JSON error: %s", err)
	}

	return tr, nil
}

// GetToken returns a new authentication and refresh token for a user
func GetToken(username, password string) (*TokenResponse, error) {
	// set form data
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", username)
	data.Set("password", password)
	return do(data)
}

// RefreshToken exchanges a refresh token for a new authentication token
func RefreshToken(refreshToken string) (*TokenResponse, error) {
	// set form data
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	return do(data)
}
