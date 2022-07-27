package constants

import "time"

const (
	AtomicRetryLimit = 5
	AtomicRetryWait  = 30 * time.Millisecond

	// KV store prefix keys
	OAuthPrefix   = "oAuth_%s"
	ProjectKey    = "%s_%s"
	ProjectPrefix = "project_list"
)
