package tado

import (
	"fmt"
	"net/http"
	"time"
)

type TerminationType string

const (
	TerminationTypeManual   TerminationType = "MANUAL"
	TerminationTypeTimer    TerminationType = "TIMER"
	TerminationTypeTadoMode TerminationType = "TADO_MODE"
)

type OverlayInput struct {
	Setting     OverlayInputSetting     `json:"setting"`
	Termination OverlayInputTermination `json:"termination"`
}

type OverlayInputSetting struct {
	Type        string                  `json:"type"`
	Power       string                  `json:"power"`
	Temperature OverlayInputTemperature `json:"temperature"`
}

type OverlayInputTermination struct {
	Type              TerminationType `json:"type"`
	DurationInSeconds int             `json:"durationInSeconds,omitempty"`
}

type OverlayInputTemperature struct {
	Celsius    float64 `json:"celsius,omitempty"`
	Fahrenheit float64 `json:"fahrenheit,omitempty"`
}

type OverlayOutput struct {
	Type    string `json:"type"`
	Setting struct {
		Type        string `json:"type"`
		Power       string `json:"power"`
		Temperature struct {
			Celsius    float64 `json:"celsius"`
			Fahrenheit float64 `json:"fahrenheit"`
		} `json:"temperature"`
	} `json:"setting"`
	Termination struct {
		Type                   string    `json:"type"`
		TypeSkillBasedApp      string    `json:"typeSkillBasedApp"`
		DurationInSeconds      int       `json:"durationInSeconds"`
		Expiry                 time.Time `json:"expiry"`
		RemainingTimeInSeconds int       `json:"remainingTimeInSeconds"`
		ProjectedExpiry        time.Time `json:"projectedExpiry"`
	} `json:"termination"`
}

// PutOverlayInput is the input for PutOverlay
type PutOverlayInput struct {
	HomeID int
	ZoneID int
	OverlayInput
}

func (poi *PutOverlayInput) method() string {
	return http.MethodPut
}

func (poi *PutOverlayInput) path() string {
	return fmt.Sprintf("/v2/homes/%d/zones/%d/overlay", poi.HomeID, poi.ZoneID)
}

func (poi *PutOverlayInput) body() interface{} {
	return poi.OverlayInput
}

// PutOverlayOutput is the output for PutOverlay
type PutOverlayOutput struct {
	OverlayOutput
}
