package constants

const (
	// Plugin API Routes
	APIPrefix                     = "/api/v1"
	WildRoute                     = "{anything:.*}"
	PathOAuthConnect              = "/oauth/connect"
	PathOAuthCallback             = "/oauth/complete"
	PathLinkedProjects            = "/project/link"
	PathGetAllLinkedProjects      = "/project/link"
	PathUnlinkProject             = "/project/unlink"
	PathUser                      = "/user"
	PathCreateTasks               = "/tasks"
	PathLinkProject               = "/link"
	PathSubscriptions             = "/subscriptions"
	PathGetSubscriptions          = "/subscriptions/{team_id:[A-Za-z0-9]+}"
	PathSubscriptionNotifications = "/notification"
	PathGetUserChannelsForTeam    = "/channels/{team_id:[A-Za-z0-9]+}"
	PathGetGitRepositories        = "/{organization:[A-Za-z0-9-/ ]+}/{project:[A-Za-z0-9-/ ]+}/repositories"
	PathGetGitRepositoryBranches  = "/{organization:[A-Za-z0-9-/ ]+}/{project:[A-Za-z0-9-/ ]+}/repositories/{repository:[A-Za-z0-9-/]+}/branches"
	PathPipelineReleaseRequest    = "/pipeline_release_request"
	PathPipelineRunRequest        = "/pipeline_run_request"

	// Azure API paths
	CreateTask                              = "/%s/%s/_apis/wit/workitems/$%s?api-version=7.1-preview.3"
	GetTask                                 = "%s/%s/_apis/wit/workitems/%s?api-version=7.1-preview.3"
	GetPullRequest                          = "%s/%s/_apis/git/pullrequests/%s?api-version=6.0"
	GetGitRepositories                      = "%s/%s/_apis/git/repositories?api-version=6.0"
	GetGitRepositoryBranches                = "%s/%s/_apis/git/repositories/%s/refs?filter=heads"
	PipelineApproveRequest                  = "%s/%s/_apis/release/approvals/%d?api-version=6.0"
	PipelineRunApproveDetails               = "/%s/%s/_apis/pipelines/approvals/%s?$expand=steps&api-version=7.0-preview.1"
	PipelineRunApproveRequest               = "%s/%s/_apis/pipelines/approvals?api-version=7.0-preview.1"
	GetProject                              = "/%s/_apis/projects/%s?api-version=7.1-preview.4"
	CreateSubscription                      = "/%s/_apis/hooks/subscriptions?api-version=6.0"
	DeleteSubscription                      = "/%s/_apis/hooks/subscriptions/%s?api-version=6.0"
	PathGetSubscriptionFilterPossibleValues = "/subscriptions/filters"
	GetSubscriptionFilterPossibleValues     = "%s/_apis/hooks/inputValuesQuery?api-version=6.0"
)
