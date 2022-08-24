package constants

import "time"

const (
	AtomicRetryLimit                     = 5
	AtomicRetryWait                      = 30 * time.Millisecond
	TTLSecondsForOAuthState        int64 = 60
	TokenExpiryTimeBufferInMinutes       = 5
	UsersPerPage                  = 100

	// KV store prefix keys
	OAuthPrefix        = "oAuth_%s"
	ProjectKey         = "%s_%s"
	ProjectPrefix      = "project_list"
	SubscriptionPrefix = "subscription_list"
	UserIDPrefix       = "oAuth"
)
