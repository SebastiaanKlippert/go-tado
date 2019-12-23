package tado

import (
	"fmt"
	"net/http"
)

type User struct {
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

// GetUsersInput is the input for GetUsers
type GetUsersInput struct {
	HomeID int
}

func (gui *GetUsersInput) method() string {
	return http.MethodGet
}

func (gui *GetUsersInput) path() string {
	return fmt.Sprintf("/v2/homes/%d/users", gui.HomeID)
}

func (gui *GetUsersInput) body() interface{} {
	return nil
}

// GetUsersOutput is the output for GetUsers
type GetUsersOutput []User
