package tado

import (
	"fmt"
	"net/http"
	"time"
)

type Device struct {
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
	InPairingMode bool `json:"inPairingMode,omitempty"`
	MountingState struct {
		Value     string    `json:"value"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"mountingState,omitempty"`
	BatteryState string `json:"batteryState,omitempty"`
}

// GetDevicesInput is the input for GetDevices
type GetDevicesInput struct {
	HomeID int
}

func (gdi *GetDevicesInput) method() string {
	return http.MethodGet
}

func (gdi *GetDevicesInput) path() string {
	return fmt.Sprintf("/v2/homes/%d/devices", gdi.HomeID)
}

func (gdi *GetDevicesInput) body() interface{} {
	return nil
}

// GetDevicesOutput is the output for GetDevices
type GetDevicesOutput []Device
