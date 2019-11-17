package tado

import (
	"net/http"
	"sync"
	"time"

	"github.com/SebastiaanKlippert/go-tado/tadoauth"
)

// Client is the main client used to communicate with the Tado API.
type Client struct {
	HTTPClient            *http.Client
	tr                    *tadoauth.TokenResponse
	accessTokenValidUntil time.Time
	mutex                 *sync.Mutex
}

// NewClient returns a new Tado client.
func NewClient(tokenResp *tadoauth.TokenResponse) *Client {
	return &Client{
		accessTokenValidUntil: time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		mutex:      new(sync.Mutex),
		tr:         tokenResp,
		HTTPClient: http.DefaultClient,
	}
}

func (c *Client) validateAccessToken() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// check if access token is valid for at least 5 more seconds
	if c.accessTokenValidUntil.After(time.Now().Add(5 * time.Second)) {
		return nil
	}
	// access token will expire soon, exchange refresh token for new access token
	tr, err := tadoauth.RefreshToken(c.tr.RefreshToken)
	if err != nil {
		return err
	}
	c.tr = tr
	c.accessTokenValidUntil = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	return nil
}

// GetMe return the users data from the API.
func (c *Client) GetMe() (*GetMeOutput, error) {
	in := new(GetMeInput)
	out := new(GetMeOutput)
	err := c.do(in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
