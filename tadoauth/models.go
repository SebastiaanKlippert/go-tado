package tadoauth

import "fmt"

// AuthenticationError is the error type returned when a Bad Request JSON message is returned from Tado servers.
// For example when the credentials provided are incorrect.
type AuthenticationError struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (ae *AuthenticationError) Error() string {
	return ae.String()
}

func (ae *AuthenticationError) String() string {
	return fmt.Sprintf("authentication error %s: %s", ae.ErrorCode, ae.ErrorDescription)
}

// TokenResponse is the response returned when a new accesstoken is returned.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	Jti          string `json:"jti"`
}
