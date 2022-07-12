package serializers

type OAuthErrorResponse struct {
	ErrorMessage     string `json:"Error"`
	ErrorDescription string `json:"ErrorDescription"`
}

type OAuthSuccessResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}
