package tado

import (
	"net/http"
	"sync"
	"time"

	"github.com/SebastiaanKlippert/go-tado/tadoauth"
)

type authClient interface {
	GetToken(username, password string) (*tadoauth.TokenResponse, error)
	RefreshToken(refreshToken string) (*tadoauth.TokenResponse, error)
}

// Client is the main client used to communicate with the Tado API.
type Client struct {
	HTTPClient *http.Client

	authClient            authClient
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

// GetMe returns the users data from the API.
func (c *Client) GetMe() (*GetMeOutput, error) {
	in := new(GetMeInput)
	out := new(GetMeOutput)
	err := c.do(in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetHome returns the Home data for a single home.
func (c *Client) GetHome(in *GetHomeInput) (*GetHomeOutput, error) {
	out := new(GetHomeOutput)
	err := c.do(in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetDevices returns the devices in a home.
func (c *Client) GetDevices(in *GetDevicesInput) (GetDevicesOutput, error) {
	out := make(GetDevicesOutput, 0)
	err := c.do(in, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetUsers returns the users in a home.
func (c *Client) GetUsers(in *GetUsersInput) (GetUsersOutput, error) {
	out := make(GetUsersOutput, 0)
	err := c.do(in, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetZones returns the zones for a single home.
func (c *Client) GetZones(in *GetZonesInput) (GetZonesOutput, error) {
	out := make(GetZonesOutput, 0)
	err := c.do(in, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetHomeState returns the presence state for a single home.
func (c *Client) GetHomeState(in *GetHomeStateInput) (*GetHomeStateOutput, error) {
	out := new(GetHomeStateOutput)
	err := c.do(in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetZoneState returns the state for a single zone within a home.
func (c *Client) GetZoneState(in *GetZoneStateInput) (*GetZoneStateOutput, error) {
	out := new(GetZoneStateOutput)
	err := c.do(in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetWeather returns the weather info for a home.
func (c *Client) GetWeather(in *GetWeatherInput) (*GetWeatherOutput, error) {
	out := new(GetWeatherOutput)
	err := c.do(in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetDayReport returns a detailed zone report of a specific date.
func (c *Client) GetDayReport(in *GetDayReportInput) (*GetDayReportOutput, error) {
	out := new(GetDayReportOutput)
	err := c.do(in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PutOverlay sets an overlay in a zone, it can be used to contol settings overruling a schema.
// For example to set the heating or hot water.
func (c *Client) PutOverlay(in *PutOverlayInput) (*PutOverlayOutput, error) {
	out := new(PutOverlayOutput)
	err := c.do(in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteOverlay deletes an overlay for a zone.
func (c *Client) DeleteOverlay(in *DeleteOverlayInput) (*DeleteOverlayOutput, error) {
	out := new(DeleteOverlayOutput)
	err := c.do(in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
