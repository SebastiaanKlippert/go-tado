package tado

import (
	"fmt"
	"net/http"
	"time"
)

// Zone contains info about a single zone
type Zone struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	DateCreated time.Time `json:"dateCreated"`
	DeviceTypes []string  `json:"deviceTypes"`
	Devices     []struct {
		DeviceType       string `json:"deviceType"`
		SerialNo         string `json:"serialNo"`
		ShortSerialNo    string `json:"shortSerialNo"`
		CurrentFwVersion string `json:"currentFwVersion"`
		ConnectionState  struct {
			Value     bool      `json:"value"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"connectionState"`
		Characteristics struct {
			Capabilities []string `json:"capabilities"`
		} `json:"characteristics"`
		MountingState struct {
			Value     string    `json:"value"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"mountingState"`
		BatteryState string   `json:"batteryState"`
		Duties       []string `json:"duties"`
	} `json:"devices"`
	ReportAvailable bool `json:"reportAvailable"`
	SupportsDazzle  bool `json:"supportsDazzle"`
	DazzleEnabled   bool `json:"dazzleEnabled"`
	DazzleMode      struct {
		Supported bool `json:"supported"`
		Enabled   bool `json:"enabled"`
	} `json:"dazzleMode"`
	OpenWindowDetection struct {
		Supported        bool `json:"supported"`
		Enabled          bool `json:"enabled"`
		TimeoutInSeconds int  `json:"timeoutInSeconds"`
	} `json:"openWindowDetection"`
}

// GetZonesInput is the input for GetZones
type GetZonesInput struct {
	HomeID int
}

func (gzi *GetZonesInput) method() string {
	return http.MethodGet
}

func (gzi *GetZonesInput) path() string {
	return fmt.Sprintf("/v2/homes/%d/zones", gzi.HomeID)
}

func (gzi *GetZonesInput) body() interface{} {
	return nil
}

// GetZonesOutput is the output for GetZones
type GetZonesOutput []Zone
