package tado

import (
	"fmt"
	"net/http"
	"time"
)

// Weather is the weather info for a home
type Weather struct {
	SolarIntensity struct {
		Type       string    `json:"type"`
		Percentage float64   `json:"percentage"`
		Timestamp  time.Time `json:"timestamp"`
	} `json:"solarIntensity"`
	OutsideTemperature struct {
		Celsius    float64   `json:"celsius"`
		Fahrenheit float64   `json:"fahrenheit"`
		Timestamp  time.Time `json:"timestamp"`
		Type       string    `json:"type"`
		Precision  struct {
			Celsius    float64 `json:"celsius"`
			Fahrenheit float64 `json:"fahrenheit"`
		} `json:"precision"`
	} `json:"outsideTemperature"`
	WeatherState struct {
		Type      string    `json:"type"`
		Value     string    `json:"value"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"weatherState"`
}

// GetWeatherInput is the input for GetWeather
type GetWeatherInput struct {
	HomeID int
}

func (gwi *GetWeatherInput) method() string {
	return http.MethodGet
}

func (gwi *GetWeatherInput) path() string {
	return fmt.Sprintf("/v2/homes/%d/weather", gwi.HomeID)
}

func (gwi *GetWeatherInput) body() interface{} {
	return nil
}

// GetWeatherOutput is the output for GetWeather
type GetWeatherOutput struct {
	Weather
}
