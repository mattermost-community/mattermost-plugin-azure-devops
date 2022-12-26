package constants

const (
	// Bot configs
	BotUsername    = "azuredevops"
	BotDisplayName = "Azure DevOps"
	BotDescription = "A bot account created by the Azure DevOps plugin."

	// Plugin configs
	PluginID               = "mattermost-plugin-azure-devops"
	ChannelID              = "channel_id"
	HeaderMattermostUserID = "Mattermost-User-ID"

	// Command configs
	CommandTriggerName = "azuredevops"
	HelpText           = "###### Mattermost Azure DevOps Plugin - Slash Command Help\n" +
		"* `/azuredevops connect` - Connect your Mattermost account to your Azure DevOps account.\n" +
		"* `/azuredevops disconnect` - Disconnect your Mattermost account from your Azure DevOps account.\n" +
		"* `/azuredevops link [projectURL]` - Link your project to a current channel.\n" +
		"* `/azuredevops boards create [title] [description]` - Create a new task for your project.\n" +
		"* `/azuredevops boards/repos/pipelines subscription add` - Add a new Boards/Repos/Pipelines subscription for your linked projects.\n" +
		"* `/azuredevops boards/repos/pipelines subscription list [me or anyone] [all_channels]` - View Boards/Repos/Pipelines subscriptions.\n" +
		"* `/azuredevops boards/repos/pipelines subscription delete [subscription id]` - Delete a Boards/Repos/Pipelines subscription"
	InvalidCommand      = "Invalid command parameters. Please use `/azuredevops help` for more information."
	CommandHelp         = "help"
	CommandConnect      = "connect"
	CommandDisconnect   = "disconnect"
	CommandLink         = "link"
	CommandBoards       = "boards"
	CommandRepos        = "repos"
	CommandPipelines    = "pipelines"
	CommandCreate       = "create"
	CommandSubscription = "subscription"
	CommandAdd          = "add"
	CommandList         = "list"
	CommandDelete       = "delete"

	// Regex to verify task link
	TaskLinkRegex = `http(s)?:\/\/dev.azure.com\/[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]*\/[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]*\/_workitems\/edit\/[1-9]+`

	// Regex to verify pull request link
	PullRequestLinkRegex = `http(s)?:\/\/dev.azure.com\/[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]*\/[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]*\/_git\/[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]*\/pullrequest\/[1-9]+`

	// Azure API Versions
	CreateTaskAPIVersion = "7.1-preview.3"
	TasksIDAPIVersion    = "5.1"
	TasksAPIVersion      = "6.0"

	// Subscription constants
	PublisherIDTFS                                      = "tfs"
	PublisherIDRM                                       = "rm"
	PublisherIDPipelines                                = "pipelines"
	ConsumerID                                          = "webHooks"
	ConsumerActionID                                    = "httpRequest"
	SubscriptionEventPullRequestCreated                 = "git.pullrequest.created"
	SubscriptionEventPullRequestUpdated                 = "git.pullrequest.updated"
	SubscriptionEventPullRequestCommented               = "ms.vss-code.git-pullrequest-comment-event"
	SubscriptionEventPullRequestMerged                  = "git.pullrequest.merged"
	SubscriptionEventCodePushed                         = "git.push"
	SubscriptionEventWorkItemCreated                    = "workitem.created"
	SubscriptionEventWorkItemUpdated                    = "workitem.updated"
	SubscriptionEventWorkItemDeleted                    = "workitem.deleted"
	SubscriptionEventWorkItemCommented                  = "workitem.commented"
	SubscriptionEventBuildCompleted                     = "build.complete"
	SubscriptionEventReleaseAbandoned                   = "ms.vss-release.release-abandoned-event"
	SubscriptionEventReleaseCreated                     = "ms.vss-release.release-created-event"
	SubscriptionEventReleaseDeploymentApprovalCompleted = "ms.vss-release.deployment-approval-completed-event"
	SubscriptionEventReleaseDeploymentEventPending      = "ms.vss-release.deployment-approval-pending-event"
	SubscriptionEventReleaseDeploymentCompleted         = "ms.vss-release.deployment-completed-event"
	SubscriptionEventReleaseDeploymentStarted           = "ms.vss-release.deployment-started-event"
	SubscriptionEventRunStageApprovalCompleted          = "ms.vss-pipelinechecks-events.approval-completed"
	SubscriptionEventRunStageStateChanged               = "ms.vss-pipelines.stage-state-changed-event"
	SubscriptionEventRunStageWaitingForApproval         = "ms.vss-pipelinechecks-events.approval-pending"
	SubscriptionEventRunStateChanged                    = "ms.vss-pipelines.run-state-changed-event"

	// Path params
	PathParamTeamID       = "team_id"
	PathParamOrganization = "organization"
	PathParamProject      = "project"
	PathParamRepository   = "repository"

	// URL query params constants
	QueryParamProject     = "project"
	QueryParamChannelID   = "channel_id"
	QueryParamCreatedBy   = "created_by"
	QueryParamServiceType = "service_type"
	QueryParamEventType   = "event_type"
	QueryParamPage        = "page"
	QueryParamPerPage     = "per_page"

	// Filters
	FilterCreatedByMe     = "me"
	FilterCreatedByAnyone = "anyone"
	FilterAllChannels     = "all_channels"
	FilterAll             = "all"
	FilterBoards          = "boards"
	FilterRepos           = "repos"

	DefaultPage         = 0
	DefaultPerPageLimit = 50

	// Authorization constants
	Bearer        = "Bearer"
	Authorization = "Authorization"

	GetTasksID  = "/%s/_apis/wit/wiql"
	GetTasks    = "/%s/_apis/wit/workitems"
	StaticFiles = "%s/plugins/%s/static/%s"

	PageQueryParam       = "$top"
	APIVersionQueryParam = "api-version"
	IDsQueryParam        = "ids"

	// Websocket events
	WSEventConnect             = "connect"
	WSEventDisconnect          = "disconnect"
	WSEventSubscriptionDeleted = "subscription_deleted"

	// Colors
	ReposIconColor  = "#d74f27"
	BoardsIconColor = "#53bba1"

	SubscriptionEventTypeDummy = "dummy"
	FileNameGitBranchIcon      = "git-branch-icon.svg"
	FileNameProjectIcon        = "project-icon.svg"
	FileNameReposIcon          = "repos-icon.svg"
	FileNameBoardsIcon         = "boards-icon.svg"
	FileNamePipelinesIcon      = "pipelines-icon.svg"
	IconColorBoards            = "#53bba1"
	IconColorRepos             = "#d74f27"
	IconColorPipelines         = "#4b68ad"

	SlackAttachmentAuthorNameRepos     = "Azure Repos"
	SlackAttachmentAuthorNameBoards    = "Azure Boards"
	SlackAttachmentAuthorNamePipelines = "Azure Pipelines"

	ServiceTypeBoards    = "boards"
	ServiceTypeRepos     = "repos"
	ServiceTypePipelines = "pipelines"

	DateTimeFormat = "Mon Jan 2 15:04:05 -0700 MST 2006"
	DateTimeLayout = "2006-01-02T15:04:05"
	TimeLayout     = "15:04:05"

	PipelineRequestIDApproved          = "approved"
	PipelineRequestIDRejected          = "rejected"
	PipelineRequestNameRun             = "run"
	PipelineRequestNameRelease         = "release"
	PipelineRequestContextApprovalID   = "approvalId"
	PipelineRequestContextOrganization = "organization"
	PipelineRequestContextProjectName  = "projectName"
	PipelineRequestContextRequestType  = "requestType"
	PipelineRequestContextRequestName  = "requestName"
	PipelineRequestContextProjectID    = "projectID"

	DialogFieldNameComment = "comment"
)

var (
	ValidSubscriptionEventsForBoards = map[string]bool{
		SubscriptionEventWorkItemCreated:   true,
		SubscriptionEventWorkItemUpdated:   true,
		SubscriptionEventWorkItemDeleted:   true,
		SubscriptionEventWorkItemCommented: true,
	}

	ValidSubscriptionEventsForRepos = map[string]bool{
		SubscriptionEventPullRequestCreated:   true,
		SubscriptionEventPullRequestMerged:    true,
		SubscriptionEventPullRequestUpdated:   true,
		SubscriptionEventPullRequestCommented: true,
		SubscriptionEventCodePushed:           true,
	}

	PipelineRequestUpdateEmoji = map[string]string{
		PipelineRequestIDApproved: "&#9989;",
		PipelineRequestIDRejected: "&#10060;",
	}
)
