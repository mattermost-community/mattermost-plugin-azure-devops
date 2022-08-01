package constants

const (
	ResponseType        = "Assertion"
	Scopes              = "vso.code vso.work_full"
	ClientAssertionType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
	GrantType           = "urn:ietf:params:oauth:grant-type:jwt-bearer"

	// URL
	BaseOauthURL = "https://app.vssps.visualstudio.com"

	// Paths
	PathAuth  = "/oauth2/authorize"
	PathToken = "/oauth2/token"
)
