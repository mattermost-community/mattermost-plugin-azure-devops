package constants

import "time"

const (
	AtomicRetryLimit              = 5
	AtomicRetryWait               = 30 * time.Millisecond
	TTLSecondsForOAuthState int64 = 60

	// KV store prefix keys
	OAuthPrefix        = "oAuth_%s"
	ProjectKey         = "%s_%s"
	ProjectPrefix      = "project_list"
	SubscriptionPrefix = "subscription_list"
)
