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
	HelpText           = "###### Mattermost Azure Devops Plugin - Slash Command Help\n" +
		"* `/azuredevops connect` - Connect your Mattermost account to your Azure Devops account.\n" +
		"* `/azuredevops disconnect` - Disconnect your Mattermost account from your Azure Devops account.\n" +
		"* `/azuredevops boards create` - Create a new task for your project.\n" +
		"* `/azuredevops link [projectURL]` - Link your project to a current channel.\n"
	InvalidCommand = "Invalid command parameters. Please use `/azuredevops help` for more information."

	// Get task link preview constants
	HTTPS              = "https:"
	HTTP               = "http:"
	AzureDevopsBaseURL = "dev.azure.com"
	Workitems          = "_workitems"
	Edit               = "edit"

	// Authorization constants
	Bearer        = "Bearer"
	Authorization = "Authorization"
)
