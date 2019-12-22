package tado

import (
	"fmt"
	"net/http"
	"time"
)

// DayReport contains the daily report info
type DayReport struct {
	ZoneType string `json:"zoneType"`
	Interval struct {
		From time.Time `json:"from"`
		To   time.Time `json:"to"`
	} `json:"interval"`
	HoursInDay   int `json:"hoursInDay"`
	MeasuredData struct {
		MeasuringDeviceConnected struct {
			TimeSeriesType string `json:"timeSeriesType"`
			ValueType      string `json:"valueType"`
			DataIntervals  []struct {
				From  time.Time `json:"from"`
				To    time.Time `json:"to"`
				Value bool      `json:"value"`
			} `json:"dataIntervals"`
		} `json:"measuringDeviceConnected"`
		InsideTemperature struct {
			TimeSeriesType string `json:"timeSeriesType"`
			ValueType      string `json:"valueType"`
			Min            struct {
				Celsius    float64 `json:"celsius"`
				Fahrenheit float64 `json:"fahrenheit"`
			} `json:"min"`
			Max struct {
				Celsius    float64 `json:"celsius"`
				Fahrenheit float64 `json:"fahrenheit"`
			} `json:"max"`
			DataPoints []struct {
				Timestamp time.Time `json:"timestamp"`
				Value     struct {
					Celsius    float64 `json:"celsius"`
					Fahrenheit float64 `json:"fahrenheit"`
				} `json:"value"`
			} `json:"dataPoints"`
		} `json:"insideTemperature"`
		Humidity struct {
			TimeSeriesType string  `json:"timeSeriesType"`
			ValueType      string  `json:"valueType"`
			PercentageUnit string  `json:"percentageUnit"`
			Min            float64 `json:"min"`
			Max            float64 `json:"max"`
			DataPoints     []struct {
				Timestamp time.Time `json:"timestamp"`
				Value     float64   `json:"value"`
			} `json:"dataPoints"`
		} `json:"humidity"`
	} `json:"measuredData"`
	Stripes struct {
		TimeSeriesType string `json:"timeSeriesType"`
		ValueType      string `json:"valueType"`
		DataIntervals  []struct {
			From  time.Time `json:"from"`
			To    time.Time `json:"to"`
			Value struct {
				StripeType string `json:"stripeType"`
				Setting    struct {
					Type        string `json:"type"`
					Power       string `json:"power"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"setting"`
			} `json:"value"`
		} `json:"dataIntervals"`
	} `json:"stripes"`
	Settings struct {
		TimeSeriesType string `json:"timeSeriesType"`
		ValueType      string `json:"valueType"`
		DataIntervals  []struct {
			From  time.Time `json:"from"`
			To    time.Time `json:"to"`
			Value struct {
				Type        string `json:"type"`
				Power       string `json:"power"`
				Temperature struct {
					Celsius    float64 `json:"celsius"`
					Fahrenheit float64 `json:"fahrenheit"`
				} `json:"temperature"`
			} `json:"value"`
		} `json:"dataIntervals"`
	} `json:"settings"`
	CallForHeat struct {
		TimeSeriesType string `json:"timeSeriesType"`
		ValueType      string `json:"valueType"`
		DataIntervals  []struct {
			From  time.Time `json:"from"`
			To    time.Time `json:"to"`
			Value string    `json:"value"`
		} `json:"dataIntervals"`
	} `json:"callForHeat"`
	Weather struct {
		Condition struct {
			TimeSeriesType string `json:"timeSeriesType"`
			ValueType      string `json:"valueType"`
			DataIntervals  []struct {
				From  time.Time `json:"from"`
				To    time.Time `json:"to"`
				Value struct {
					State       string `json:"state"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"value"`
			} `json:"dataIntervals"`
		} `json:"condition"`
		Sunny struct {
			TimeSeriesType string `json:"timeSeriesType"`
			ValueType      string `json:"valueType"`
			DataIntervals  []struct {
				From  time.Time `json:"from"`
				To    time.Time `json:"to"`
				Value bool      `json:"value"`
			} `json:"dataIntervals"`
		} `json:"sunny"`
		Slots struct {
			TimeSeriesType string `json:"timeSeriesType"`
			ValueType      string `json:"valueType"`
			Slots          struct {
				Zero400 struct {
					State       string `json:"state"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"04:00"`
				Zero800 struct {
					State       string `json:"state"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"08:00"`
				One200 struct {
					State       string `json:"state"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"12:00"`
				One600 struct {
					State       string `json:"state"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"16:00"`
				Two000 struct {
					State       string `json:"state"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"20:00"`
			} `json:"slots"`
		} `json:"slots"`
	} `json:"weather"`
}

// GetDayReportInput is the input for GetDayReport
type GetDayReportInput struct {
	HomeID int
	ZoneID int
	Date   time.Time
}

func (gdri *GetDayReportInput) method() string {
	return http.MethodGet
}

func (gdri *GetDayReportInput) path() string {
	return fmt.Sprintf("/v2/homes/%d/zones/%d/dayReport?date=%s", gdri.HomeID, gdri.ZoneID, gdri.Date.Format("2006-01-02"))
}

func (gdri *GetDayReportInput) body() interface{} {
	return nil
}

// GetDayReportOutput is the output for GetDayReport
type GetDayReportOutput struct {
	DayReport
}
