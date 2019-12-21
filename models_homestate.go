package tado

import (
	"fmt"
	"net/http"
)

const (
	HomeStateHome = "HOME" // HomeStateHome is the value of presence used when home
	HomeStateAway = "AWAY" // HomeStateAway is the value of presence used when away
)

// HomeState is the state of a home
type HomeState struct {
	Presence string `json:"presence"`
}

// GetHomeStateInput is the input for GetHomeState
type GetHomeStateInput struct {
	HomeID int
}

func (ghsi *GetHomeStateInput) method() string {
	return http.MethodGet
}

func (ghsi *GetHomeStateInput) path() string {
	return fmt.Sprintf("/v2/homes/%d/state", ghsi.HomeID)
}

func (ghsi *GetHomeStateInput) body() interface{} {
	return nil
}

// GetHomeStateOutput is the output for GetHomeState
type GetHomeStateOutput struct {
	HomeState
}
