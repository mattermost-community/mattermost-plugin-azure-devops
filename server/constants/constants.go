package constants

const (
	// Bot configs
	BotUsername    = "azuredevops"
	BotDisplayName = "Azure Devops"
	BotDescription = "A bot account created by the Azure Devops plugin."

	// Plugin configs
	PluginID               = "mattermost-plugin-azure-devops"
	ChannelID              = "channel_id"
	HeaderMattermostUserID = "Mattermost-User-ID"
	// TODO: Change later according to the needs.
	HeaderMattermostUserIDAPI = "User-ID"

	// Command configs
	CommandTriggerName = "azuredevops"
	HelpText           = "###### Mattermost Azure Devops Plugin - Slash Command Help\n"
	InvalidCommand     = "Invalid command parameters. Please use `/azuredevops help` for more information."

	// Azure API Routes
	CreateTask = "/%s/%s/_apis/wit/workitems/$%s?api-version=" + CreateTaskAPIVersion

	// Azure API Versions
	CreateTaskAPIVersion = "7.1-preview.3"

	// Authorization constants
	Bearer        = "Bearer"
	Authorization = "Authorization"

	GetTasksID = "/%s/_apis/wit/wiql"
	GetTasks   = "/%s/_apis/wit/workitems"

	TasksIDAPIVersion = "5.1"
	TasksAPIVersion   = "6.0"

	PageQueryParam       = "$top"
	APIVersionQueryParam = "api-version"
	IDsQueryParam        = "ids"
)
