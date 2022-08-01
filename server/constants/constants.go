package constants

const (
	// Bot configs
	BotUsername    = "azuredevops"
	BotDisplayName = "Azure Devops"
	BotDescription = "A bot account created by the Azure Devops plugin."

	// Plugin configs
	PluginID                  = "mattermost-plugin-azure-devops"
	HeaderMattermostUserID    = "Mattermost-User-ID"
	HeaderMattermostUserIDAPI = "User-ID"
	ChannelID                 = "channel_id"

	// Command configs
	CommandTriggerName = "azuredevops"
	HelpText           = "###### Mattermost Azure Devops Plugin - Slash Command Help\n"
	InvalidCommand     = "Invalid command parameters. Please use `/azuredevops help` for more information."

	// Get task link preview constants
	HTTPS              = "https:"
	HTTP               = "http:"
	AzureDevopsBaseURL = "dev.azure.com"
	Workitems          = "_workitems"
	Edit               = "edit"

	// Authorization constants
	Bearer        = "Bearer %s"
	Authorization = "Authorization"
)
