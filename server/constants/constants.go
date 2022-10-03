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
		"* `/azuredevops boards subscription add` - Add a new Boards subscription for your linked projects.\n" +
		"* `/azuredevops boards subscription list [me or anyone] [all-channels]` - View Boards subscriptions.\n" +
		"* `/azuredevops boards subscription delete [subscription id]` - Unsubscribe a Boards subscription"
	InvalidCommand      = "Invalid command parameters. Please use `/azuredevops help` for more information."
	CommandHelp         = "help"
	CommandConnect      = "connect"
	CommandDisconnect   = "disconnect"
	CommandLink         = "link"
	CommandBoards       = "boards"
	CommandCreate       = "create"
	CommandSubscription = "subscription"
	CommandAdd          = "add"
	CommandList         = "list"
	CommandDelete       = "delete"

	// Get task link preview constants
	HTTPS              = "https:"
	HTTP               = "http:"
	AzureDevopsBaseURL = "dev.azure.com"
	Workitems          = "_workitems"
	Edit               = "edit"

	// Azure API Versions
	CreateTaskAPIVersion = "7.1-preview.3"
	TasksIDAPIVersion    = "5.1"
	TasksAPIVersion      = "6.0"

	// Subscription constants
	PublisherID      = "tfs"
	ConsumerID       = "webHooks"
	ConsumerActionID = "httpRequest"
	Create           = "create"
	Update           = "update"
	Delete           = "delete"

	// Path params
	PathParamTeamID = "team_id"

	// URL query params constants
	QueryParamProject   = "project"
	QueryParamChannelID = "channel_id"
	QueryParamCreatedBy = "created_by"
	QueryParamPage      = "page"
	QueryParamPerPage   = "per_page"

	// Filters
	FilterCreatedByMe     = "me"
	FilterCreatedByAnyone = "anyone"
	FilterAllChannels     = "all-channels"
	FilterCurrentChannel  = "current-channel"

	DefaultPage         = 0
	DefaultPerPageLimit = 50

	// Authorization constants
	Bearer        = "Bearer"
	Authorization = "Authorization"

	GetTasksID = "/%s/_apis/wit/wiql"
	GetTasks   = "/%s/_apis/wit/workitems"

	PageQueryParam       = "$top"
	APIVersionQueryParam = "api-version"
	IDsQueryParam        = "ids"

	// Websocket events
	WSEventConnect             = "connect"
	WSEventDisconnect          = "disconnect"
	WSEventSubscriptionDeleted = "subscription_deleted"
)
