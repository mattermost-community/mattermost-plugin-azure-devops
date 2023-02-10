package serializers

type GenerateTokenPayload struct {
	ClientAssertionType string `json:"client_assertion_type"`
	ClientAssertion     string `json:"client_assertion"`
	GrantType           string `json:"grant_type"`
	Assertion           string `json:"assertion"`
	RedirectURI         string `json:"redirect_uri"`
}

type OAuthErrorResponse struct {
	ErrorMessage     string `json:"Error"`
	ErrorDescription string `json:"ErrorDescription"`
}

type OAuthSuccessResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

type ConnectedResponse struct {
	IsConnected bool `json:"connected"`
}

type UserProfile struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"emailAddress"`
}
