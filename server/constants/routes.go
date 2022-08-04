package constants

const (
	// Plugin API Routes
	APIPrefix                     = "/api/v1"
	WildRoute                     = "{anything:.*}"
	PathOAuthConnect              = "/oauth/connect"
	PathOAuthCallback             = "/oauth/complete"
	PathGetAllLinkedProjects      = "/project/link"
	PathUnlinkProject             = "/project/unlink"
	PathUser                      = "/user"
	PathCreateTasks               = "/tasks"
	PathLinkProject               = "/link"
	PathGetSubscriptions          = "/subscriptions"
	PathCreateSubscriptions       = "/subscriptions"
	PathNotificationSubscriptions = "/notification"

	// Azure API paths
	CreateTask         = "/%s/%s/_apis/wit/workitems/$%s?api-version=7.1-preview.3"
	GetTask            = "%s/_apis/wit/workitems/%s?api-version=7.1-preview.3"
	GetProject         = "/%s/_apis/projects/%s?api-version=7.1-preview.4"
	CreateSubscription = "/%s/_apis/hooks/subscriptions?api-version=6.0"
)
