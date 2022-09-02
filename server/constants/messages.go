package constants

const (
	// TODO: all these messages are to be verified from Mike at the end
	GenericErrorMessage  = "Something went wrong, please try again later."
	ConnectAccount       = "[Click here to link your Azure DevOps account](%s%s)"
	ConnectAccountFirst  = "You do not have any Azure DevOps account connected. Kindly connect your account first."
	UserConnected        = "Your Azure DevOps account is succesfully connected!"
	UserAlreadyConnected = "Your Azure DevOps account is already connected."
	UserDisconnected     = "Your Azure DevOps account is now disconnected."
	CreatedTask          = "Link for newly created task: %s"
	TaskTitle            = "[%s #%d: %s](%s)"
	TaskPreviewMessage   = "State: %s\nAssigned To: %s\nDescription: %s"
	AlreadyLinkedProject = "This project is already linked."
	NoProjectLinked      = "No project is linked, please link a project."

	// Validations
	OrganizationRequired = "organization is required"
	ProjectRequired      = "project is required"
	TaskTypeRequired     = "task type is required"
	TaskTitleRequired    = "task title is required"
	EventTypeRequired    = "event type is required"
	ChannelNameRequired  = "channel name is required"
	ChannelIDRequired    = "channel ID is required"
)

const (
	// Error messages
	Error                           = "error"
	NotAuthorized                   = "not authorized"
	UnableToDisconnectUser          = "Unable to disconnect user"
	UnableToCheckIfAlreadyConnected = "Unable to check if user account is already connected"
	UnableToStoreOauthState         = "Unable to store oAuth state for the userID %s"
	AuthAttemptExpired              = "Authentication attempt expired, please try again"
	InvalidAuthState                = "Invalid oauth state, please try again"
	GetProjectListError             = "Error in getting project list"
	ErrorFetchProjectList           = "Error in fetching project list"
	ErrorDecodingBody               = "Error in decoding body"
	ErrorCreateTask                 = "Error in creating task"
	ErrorLinkProject                = "Error in linking the project"
	ErrorLoadingDataFromKVStore     = "Error in loading data from KV store"
	ProjectNotFound                 = "Requested project does not exist"
	ErrorUnlinkProject              = "Error in unlinking the project"
	FetchSubscriptionListError      = "Error in fetching subscription list"
	CreateSubscriptionError         = "Error in creating subscription"
	ProjectNotLinked                = "Requested project is not linked"
	GetSubscriptionListError        = "Error in getting subscription list"
	SubscriptionAlreadyPresent      = "Subscription is already present"
	InvalidChannelID                = "Invalid channel ID"
)
