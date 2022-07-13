package constants

const (
	// Bot configs
	BotUsername    = "azuredevops"
	BotDisplayName = "Azure Devops Plugin"
	BotDescription = "A bot account created by the Azure Devops plugin."

	// Plugin configs
	PluginID               = "mattermost-plugin-azure-devops"
	HeaderMattermostUserID = "Mattermost-User-ID"

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
	GetProject = "/%s/_apis/projects?$top=%d&api-version=7.1-preview.4"
	GetTasksID = "/%s/_apis/wit/wiql?$top=%d&api-version=5.1"
	GetTasks   = "/%s/_apis/wit/workitems?ids=%s&api-version=6.0"

	// Max tasks and projects per page
	MaxProjectsPerPage = 10
	MaxTasksPerPage    = 10

	// Generic messages
	// TODO: all these messages are to be verified from Mike at the end
	ConnectAccount       = "[Click here to link your Azure DevOps account](%s%s?channel_id=%s)"
	ConnectAccountFirst  = "You do not have any Azure Devops account connected, kindly link the account first"
	UserConnected        = "Your Azure Devops account is succesfully connected!"
	UserAlreadyConnected = "Your Azure Devops account is already connected"
	UserDisconnected     = "Your Azure Devops account is now disconnected"

	// Error messages
	GenericErrorMessage  = "something went wrong, please try again later"
	NotAuthorized        = "not authorized"
	InvalidPageNumber    = "invalid page number"
	InvalidLimitNumber   = "invalid limit"
	OrganizationRequired = "organization is required"
	ProjectRequired      = "project is required"
	InvalidStatus        = "invalid status"
	InvalidAssignedTo    = "you can only see tasks assigned to yourself"
	NoResultPresent      = "no results are present"
)
