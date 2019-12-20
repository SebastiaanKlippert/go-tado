package tado

import (
	"net/http"
	"sync"
	"time"

	"github.com/SebastiaanKlippert/go-tado/tadoauth"
)

// Client is the main client used to communicate with the Tado API.
type Client struct {
	HTTPClient *http.Client

	authClient            *tadoauth.Client
	baseURL               string
	username, password    string
	tr                    *tadoauth.TokenResponse
	accessTokenValidUntil time.Time
	mutex                 *sync.Mutex
}

// NewClient returns a new Tado client.
func NewClient(username, password string) *Client {
	return &Client{
		username:              username,
		password:              password,
		accessTokenValidUntil: time.Time{},
		baseURL:               defaultBaseURL,
		mutex:                 new(sync.Mutex),
		authClient:            tadoauth.NewClient(),
		HTTPClient:            http.DefaultClient,
	}
}

func (c *Client) validateAccessToken() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// check if access token is valid for at least 5 more seconds
	if c.accessTokenValidUntil.After(time.Now().Add(5 * time.Second)) {
		return nil
	}
	// access token is expired, or will expire soon, get a new one
	var err error
	if c.tr != nil && c.tr.RefreshToken != "" {
		// exchange refresh token for new access token
		c.tr, err = c.authClient.RefreshToken(c.tr.RefreshToken)
	} else {
		// get new token based on username and password
		c.tr, err = c.authClient.GetToken(c.username, c.password)
	}
	if err != nil {
		return err
	}
	c.accessTokenValidUntil = time.Now().Add(time.Duration(c.tr.ExpiresIn) * time.Second)
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

// GetHome return the Home data for a single home ID.
func (c *Client) GetHome(in *GetHomeInput) (*GetHomeOutput, error) {
	out := new(GetHomeOutput)
	err := c.do(in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetZones return the zones for a single home ID.
func (c *Client) GetZones(in *GetZonesInput) (GetZonesOutput, error) {
	out := make(GetZonesOutput, 0)
	err := c.do(in, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
