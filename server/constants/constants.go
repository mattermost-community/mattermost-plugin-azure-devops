package constants

const (
	// Bot configs
	BotUsername    = "azuredevops"
	BotDisplayName = "Azure Devops Plugin"
	BotDescription = "A bot account created by the Azure Devops plugin."

	// Command configs
	CommandTriggerName = "azuredevops"
	HelpText           = "###### Mattermost Azure Devops Plugin - Slash Command Help\n"
	InvalidCommand     = "Invalid command parameters. Please use `/azuredevops help` for more information."

	// Plugin API Routes
	APIPrefix = "/api/v1"
	WildRoute = "{anything:.*}"

	// Error messages
	GenericErrorMessage = "something went wrong, please try again later"
)
