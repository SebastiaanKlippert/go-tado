package tado

import "net/http"

// Me contains the users data.
type Me struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Enabled  bool   `json:"enabled"`
	ID       string `json:"id"`
	HomeID   int    `json:"homeId"`
	Locale   string `json:"locale"`
	Type     string `json:"type"`
}

// GetMeInput is the input for GetMe
type GetMeInput struct{}

func (gmi *GetMeInput) method() string {
	return http.MethodGet
}

func (gmi *GetMeInput) path() string {
	return "/v1/me"
}

// GetMeOutput is the output for GetMe
type GetMeOutput struct {
	Me
}
