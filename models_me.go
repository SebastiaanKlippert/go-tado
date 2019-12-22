package tado

import "net/http"

// Me contains the users data
type Me struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	ID       string `json:"id"`
	Homes    []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"homes"`
	Locale        string `json:"locale"`
	MobileDevices []struct {
		Name     string `json:"name"`
		ID       int    `json:"id"`
		Settings struct {
			GeoTrackingEnabled bool `json:"geoTrackingEnabled"`
		} `json:"settings"`
		DeviceMetadata struct {
			Platform  string `json:"platform"`
			OsVersion string `json:"osVersion"`
			Model     string `json:"model"`
			Locale    string `json:"locale"`
		} `json:"deviceMetadata"`
	} `json:"mobileDevices"`
}

// GetMeInput is the input for GetMe
type GetMeInput struct{}

func (gmi *GetMeInput) method() string {
	return http.MethodGet
}

func (gmi *GetMeInput) path() string {
	return "/v2/me"
}

func (gmi *GetMeInput) body() interface{} {
	return nil
}

// GetMeOutput is the output for GetMe
type GetMeOutput struct {
	Me
}
