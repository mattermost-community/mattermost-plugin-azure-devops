package constants

const (
	// Bot configs
	BotUsername    = "azuredevops"
	BotDisplayName = "Azure Devops Plugin"
	BotDescription = "A bot account created by the Azure Devops plugin."

	// Plugin configs
	PluginID = "mattermost-plugin-azure-devops"
	// TODO: Change later according to the needs.
	HeaderMattermostUserID = "User-ID"

	// Command configs
	CommandTriggerName = "azuredevops"
	HelpText           = "###### Mattermost Azure Devops Plugin - Slash Command Help\n"
	InvalidCommand     = "Invalid command parameters. Please use `/azuredevops help` for more information."

	// Plugin API Routes
	APIPrefix         = "/api/v1"
	WildRoute         = "{anything:.*}"
	PathOAuthConnect  = "/connect"
	PathOAuthCallback = "/callback"

	// Azure API Routes
	// TODO: WIP.
	// GetProjects = "/%s/_apis/projects"
	// GetTasksID = "/%s/_apis/wit/wiql"
	// GetTasks   = "/%s/_apis/wit/workitems"
	CreateTask = "/%s/%s/_apis/wit/workitems/$%s"

	// Azure API versions
	// TODO: WIP.
	// ProjectAPIVersion = "7.1-preview.4"
	// TasksIDAPIVersion = "5.1"
	// TasksAPIVersion   = "6.0"
	CreateTaskAPIVersion = "7.1-preview.3"

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
	// IDsQueryParam        = "ids"
	APIVersionQueryParam = "api-version"

	// Generic messages
	// TODO: all these messages are to be verified from Mike at the end
	ConnectAccount       = "[Click here to link your Azure DevOps account](%s%s?channel_id=%s)"
	ConnectAccountFirst  = "You do not have any Azure Devops account connected, kindly link the account first"
	UserConnected        = "Your Azure Devops account is succesfully connected!"
	UserAlreadyConnected = "Your Azure Devops account is already connected"
	UserDisconnected     = "Your Azure Devops account is now disconnected"

	// Error messages
	Error               = "error"
	GenericErrorMessage = "something went wrong, please try again later"
	NotAuthorized       = "not authorized"
	// TODO: WIP.
	// InvalidPageNumber    = "invalid page number"
	// InvalidStatus        = "invalid status"
	// InvalidAssignedTo    = "you can only see tasks assigned to yourself"
	// NoResultPresent      = "no results are present"
	OrganizationRequired = "organization is required"
	ProjectRequired      = "project is required"
	TaskTypeRequired     = "task type is required"
	TaskTitleRequired    = "task title is required"
)
