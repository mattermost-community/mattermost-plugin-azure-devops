package constants

const (
	// Bot configs
	BotUsername    = "azuredevops"
	BotDisplayName = "Azure Devops"
	BotDescription = "A bot account created by the Azure Devops plugin."

	// Plugin configs
	PluginID               = "mattermost-plugin-azure-devops"
	HeaderMattermostUserID = "Mattermost-User-ID"
	// TODO: Change later according to the needs.
	HeaderMattermostUserIDAPI = "User-ID"
	ChannelID                 = "channel_id"

	// Command configs
	CommandTriggerName = "azuredevops"
	HelpText           = "###### Mattermost Azure Devops Plugin - Slash Command Help\n"
	InvalidCommand     = "Invalid command parameters. Please use `/azuredevops help` for more information."

	// Azure API Routes
	// TODO: WIP.
	// GetProjects = "/%s/_apis/projects"
	GetTasksID = "/%s/_apis/wit/wiql"
	GetTasks   = "/%s/_apis/wit/workitems"
	CreateTask = "/%s/%s/_apis/wit/workitems/$%s?api-version=7.1-preview.3"
	GetTask    = "%s/_apis/wit/workitems/%s?api-version=7.1-preview.3"

	// Azure API versions
	// TODO: WIP.
	// ProjectAPIVersion = "7.1-preview.4"
	// TasksIDAPIVersion    = "5.1"
	// TasksAPIVersion      = "6.0"
	CreateTaskAPIVersion = "7.1-preview.3"

	// Get task link preview constants
	HTTPS              = "https:"
	HTTP               = "http:"
	AzureDevopsBaseURL = "dev.azure.com"
	Workitems          = "_workitems"
	Edit               = "edit"

	// Authorization constants
	Bearer        = "Bearer %s"
	Authorization = "Authorization"

	// Limits
	// TODO: WIP.
	// ProjectLimit = 10
	// TaskLimit = 10

	// URL filters
	Organization = "organization"
	Project      = "project"
	Status       = "status"
	AssignedTo   = "assigned_to"
	Page         = "page"

	// Tasks status
	Doing = "doing"
	Todo  = "to-do"
	Done  = "done"

	// Query params constants
	// TODO: WIP.
	// PageQueryParam       = "$top"
	// IDsQueryParam        = "ids"

	// Generic messages
	// TODO: all these messages are to be verified from Mike at the end
	ConnectAccount       = "[Click here to link your Azure DevOps account](%s%s?channel_id=%s)"
	ConnectAccountFirst  = "You do not have any Azure Devops account connected, kindly link the account first"
	UserConnected        = "Your Azure Devops account is succesfully connected!"
	UserAlreadyConnected = "Your Azure Devops account is already connected"
	UserDisconnected     = "Your Azure Devops account is now disconnected"
	CreatedTask          = "Link for newly created task: %s"
	TaskTitle            = "[%s #%d: %s](%s)"
	TaskPreviewMessage   = "State: %s\nAssigned To: %s\nDescription: %s"

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
