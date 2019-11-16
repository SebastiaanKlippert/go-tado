package tado_auth

import "fmt"

type AthenticationError struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (ae *AthenticationError) Error() string {
	return ae.String()
}

func (ae *AthenticationError) String() string {
	return fmt.Sprintf("authentication error %s: %s", ae.ErrorCode, ae.ErrorDescription)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	Jti          string `json:"jti"`
}
