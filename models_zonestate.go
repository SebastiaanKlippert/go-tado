package tado

import (
	"fmt"
	"net/http"
	"time"
)

// ZoneState is the state of a single zone
type ZoneState struct {
	TadoMode                       string      `json:"tadoMode"`
	GeolocationOverride            bool        `json:"geolocationOverride"`
	GeolocationOverrideDisableTime interface{} `json:"geolocationOverrideDisableTime"` // TODO
	Preparation                    interface{} `json:"preparation"`                    // TODO
	Setting                        struct {
		Type        string `json:"type"`
		Power       string `json:"power"`
		Temperature struct {
			Celsius    float64 `json:"celsius"`
			Fahrenheit float64 `json:"fahrenheit"`
		} `json:"temperature"`
	} `json:"setting"`
	OverlayType        interface{} `json:"overlayType"`        // TODO
	Overlay            interface{} `json:"overlay"`            // TODO
	OpenWindow         interface{} `json:"openWindow"`         // TODO
	NextScheduleChange interface{} `json:"nextScheduleChange"` // TODO
	NextTimeBlock      interface{} `json:"nextTimeBlock"`      // TODO
	Link               struct {
		State string `json:"state"`
	} `json:"link"`
	ActivityDataPoints struct {
		HeatingPower struct {
			Type       string    `json:"type"`
			Percentage float64   `json:"percentage"`
			Timestamp  time.Time `json:"timestamp"`
		} `json:"heatingPower"`
	} `json:"activityDataPoints"`
	SensorDataPoints struct {
		InsideTemperature struct {
			Celsius    float64   `json:"celsius"`
			Fahrenheit float64   `json:"fahrenheit"`
			Timestamp  time.Time `json:"timestamp"`
			Type       string    `json:"type"`
			Precision  struct {
				Celsius    float64 `json:"celsius"`
				Fahrenheit float64 `json:"fahrenheit"`
			} `json:"precision"`
		} `json:"insideTemperature"`
		Humidity struct {
			Type       string    `json:"type"`
			Percentage float64   `json:"percentage"`
			Timestamp  time.Time `json:"timestamp"`
		} `json:"humidity"`
	} `json:"sensorDataPoints"`
}

// GetZoneStateInput is the input for GetZoneState
type GetZoneStateInput struct {
	HomeID int
	ZoneID int
}

func (gzsi *GetZoneStateInput) method() string {
	return http.MethodGet
}

func (gzsi *GetZoneStateInput) path() string {
	return fmt.Sprintf("/v2/homes/%d/zones/%d/state", gzsi.HomeID, gzsi.ZoneID)
}

func (gzsi *GetZoneStateInput) body() interface{} {
	return nil
}

// GetZoneStateOutput is the output for GetZoneState
type GetZoneStateOutput struct {
	ZoneState
}
