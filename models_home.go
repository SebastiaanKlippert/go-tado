package tado

import (
	"fmt"
	"net/http"
	"time"
)

type Home struct {
	ID                         int           `json:"id"`
	Name                       string        `json:"name"`
	DateTimeZone               string        `json:"dateTimeZone"`
	DateCreated                time.Time     `json:"dateCreated"`
	TemperatureUnit            string        `json:"temperatureUnit"`
	InstallationCompleted      bool          `json:"installationCompleted"`
	Partner                    string        `json:"partner"`
	SimpleSmartScheduleEnabled bool          `json:"simpleSmartScheduleEnabled"`
	AwayRadiusInMeters         float64       `json:"awayRadiusInMeters"`
	UsePreSkillsApps           bool          `json:"usePreSkillsApps"`
	Skills                     []interface{} `json:"skills"` //TODO define
	ChristmasModeEnabled       bool          `json:"christmasModeEnabled"`
	ShowAutoAssistReminders    bool          `json:"showAutoAssistReminders"`
	ContactDetails             struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"contactDetails"`
	Address struct {
		AddressLine1 string `json:"addressLine1"`
		AddressLine2 string `json:"addressLine2"`
		ZipCode      string `json:"zipCode"`
		City         string `json:"city"`
		State        string `json:"state"`
		Country      string `json:"country"`
	} `json:"address"`
	Geolocation struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"geolocation"`
	ConsentGrantSkippable bool `json:"consentGrantSkippable"`
}

// GetHomeInput is the input for GetMe
type GetHomeInput struct {
	HomeID int
}

func (ghi *GetHomeInput) method() string {
	return http.MethodGet
}

func (ghi *GetHomeInput) path() string {
	return fmt.Sprintf("/v2/homes/%d", ghi.HomeID)
}

// GetHomeOutput is the output for GetMe
type GetHomeOutput struct {
	Home
}
