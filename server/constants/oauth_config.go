package constants

const (
	ResponseType        = "Assertion"
	Scopes              = "vso.build_execute vso.code_full vso.release_manage vso.work_full"
	ClientAssertionType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
	GrantType           = "urn:ietf:params:oauth:grant-type:jwt-bearer"
	GrantTypeRefresh    = "refresh_token"

	// URL
	BaseOauthURL = "https://app.vssps.visualstudio.com"

	// Paths
	PathAuth = "/oauth2/authorize"
	// #nosec G101 -- This is a false positive
	PathToken       = "/oauth2/token"
	PathUserProfile = "/_apis/profile/profiles/%s"

	CurrentAzureDevopsUserProfileID = "me"
)
