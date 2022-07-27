package constants

const (
	// Bot configs
	BotUsername    = "azuredevops"
	BotDisplayName = "Azure Devops"
	BotDescription = "A bot account created by the Azure Devops plugin."

	// Plugin configs
	PluginID  = "mattermost-plugin-azure-devops"
	ChannelID = "channel_id"
	// TODO: Change later according to the needs.
	HeaderMattermostUserID = "User-ID"

	// Command configs
	CommandTriggerName = "azuredevops"
	HelpText           = "###### Mattermost Azure Devops Plugin - Slash Command Help\n"
	InvalidCommand     = "Invalid command parameters. Please use `/azuredevops help` for more information."

	// Azure API Routes
	// TODO: WIP.
	// GetProjects = "/%s/_apis/projects"
	// GetTasksID = "/%s/_apis/wit/wiql"
	// GetTasks   = "/%s/_apis/wit/workitems"

	// Azure API versions
	// TODO: WIP.
	// ProjectAPIVersion = "7.1-preview.4"
	// TasksIDAPIVersion = "5.1"
	// TasksAPIVersion   = "6.0"

	// Authorization constants
	Bearer        = "Bearer %s"
	Authorization = "Authorization"

	// Limits
	// TODO: WIP.
	// ProjectLimit = 10
	// TaskLimit = 10

	// TODO: WIP.
	// URL filters
	// Organization = "organization"
	// Project      = "project"
	// Status       = "status"
	// AssignedTo   = "assigned_to"
	// Page         = "page"

	// TODO: WIP.
	// Tasks status
	// Doing = "doing"
	// Todo  = "to-do"
	// Done  = "done"

	// TODO: WIP.
	// Query params constants
	// PageQueryParam       = "$top"
	// APIVersionQueryParam = "api-version"
	// IDsQueryParam        = "ids"
)
